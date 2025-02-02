import json
import logging
import os
from datetime import datetime, timedelta
from typing import Any, Dict, List, Optional

import pendulum
import requests
from airflow.configuration import conf
from airflow.exceptions import AirflowException
from airflow.hooks.base import BaseHook
from airflow.kubernetes import kube_client
from airflow.models import (XCOM_RETURN_KEY, DagModel,
                            DagRun, Variable, XCom)
from airflow.providers.cncf.kubernetes.operators.kubernetes_pod import KubernetesPodOperator
from airflow.providers.cncf.kubernetes.utils import pod_launcher
from airflow.providers.slack.operators.slack import SlackAPIPostOperator
from airflow.sensors.base_sensor_operator import BaseSensorOperator
from airflow.utils.db import provide_session
from airflow.utils.decorators import apply_defaults
from airflow.utils.state import State
from croniter import croniter

log = logging.getLogger(__name__)
log.setLevel(logging.INFO)

# UTC time zone as a tzinfo instance.
utc = pendulum.timezone('UTC')

TIMESTAMP_FORMAT = "%Y-%m-%dT%H:%M:%SZ"
TIMESTAMP_MS_FORMAT = "%Y-%m-%dT%H:%M:%S.%fZ"


def lookup_non_standard_cron_expression(expr: str) -> str:
    expr_mapping = {
        '@yearly': '0 0 1 1 *',
        '@annually': '0 0 1 1 *',
        '@monthly': '0 0 1 * *',
        '@weekly': '0 0 * * 0',
        '@daily': '0 0 * * *',
        '@midnight': '0 0 * * *',
        '@hourly': '0 * * * *',
    }
    try:
        return expr_mapping[expr]
    except KeyError:
        return expr


class SuperKubernetesPodOperator(KubernetesPodOperator):
    """
    ** SAME AS KubernetesPodOperator: Execute a task in a Kubernetes Pod **
    Wrapper to push xcom as a return value key even if container completes with non success status

    .. note: keep this up to date if there is any change in KubernetesPodOperator execute method
    """
    template_fields = ('image', 'cmds', 'arguments', 'env_vars', 'config_file', 'pod_template_file')

    @apply_defaults
    def __init__(self,
                 *args,
                 **kwargs):
        super(SuperKubernetesPodOperator, self).__init__(*args, **kwargs)

        self.do_xcom_push = kwargs.get('do_xcom_push')
        self.namespace = kwargs.get('namespace')
        self.in_cluster = kwargs.get('in_cluster')
        self.cluster_context = kwargs.get('cluster_context')
        self.reattach_on_restart = kwargs.get('reattach_on_restart')
        self.config_file = kwargs.get('config_file')

    def execute(self, context):
        log.info('Task image version: %s', self.image)
        try:
            if self.in_cluster is not None:
                client = kube_client.get_kube_client(in_cluster=self.in_cluster,
                                                     cluster_context=self.cluster_context,
                                                     config_file=self.config_file)
            else:
                client = kube_client.get_kube_client(cluster_context=self.cluster_context,
                                                     config_file=self.config_file)

            self.pod = self.create_pod_request_obj()
            self.namespace = self.pod.metadata.namespace
            self.client = client

            # Add combination of labels to uniquely identify a running pod
            labels = self.create_labels_for_pod(context)

            label_selector = self._get_pod_identifying_label_string(labels)

            pod_list = client.list_namespaced_pod(self.namespace, label_selector=label_selector)

            if len(pod_list.items) > 1 and self.reattach_on_restart:
                raise AirflowException(
                    'More than one pod running with labels: '
                    '{label_selector}'.format(label_selector=label_selector))

            launcher = pod_launcher.PodLauncher(kube_client=client, extract_xcom=self.do_xcom_push)

            if len(pod_list.items) == 1:
                try_numbers_match = self._try_numbers_match(context, pod_list.items[0])
                final_state, result = self.handle_pod_overlap(labels, try_numbers_match, launcher, pod_list.items[0])
            else:
                final_state, _, result = self.create_new_pod_for_operator(labels, launcher)

            if final_state != State.SUCCESS:
                # push xcom value even if pod fails
                context.get('task_instance').xcom_push(key=XCOM_RETURN_KEY, value=result)
                raise AirflowException(
                    'Pod returned a failure: {state}'.format(state=final_state))
            return result
        except AirflowException as ex:
            raise AirflowException('Pod Launching failed: {error}'.format(error=ex))

class OptimusAPIClient:
    def __init__(self, optimus_host):
        self.host = self._add_connection_adapter_if_absent(optimus_host)

    def _add_connection_adapter_if_absent(self, host):
        if host.startswith("http://") or host.startswith("https://"):
            return host
        return "http://" + host

    def get_job_run(self, optimus_project: str, optimus_job: str, startDate: str, endDate:str ) -> dict:
        url = '{optimus_host}/api/v1beta1/project/{optimus_project}/job/{optimus_job}/run'.format(
            optimus_host=self.host,
            optimus_project=optimus_project,
            optimus_job=optimus_job,
        )
        response = requests.get(url, params = { 'start_date': startDate,'end_date': endDate})
        self._raise_error_if_request_failed(response)
        return response.json()

    def get_task_window(self, scheduled_at: str, window_size: str, window_offset: str,
                        window_truncate_upto: str) -> dict:
        url = '{optimus_host}/api/v1beta1/window?scheduledAt={scheduled_at}&size={window_size}&offset={window_offset}&truncate_to={window_truncate_upto}'.format(
            optimus_host=self.host,
            scheduled_at=scheduled_at,
            window_size=window_size,
            window_offset=window_offset,
            window_truncate_upto=window_truncate_upto,
        )
        response = requests.get(url)
        self._raise_error_if_request_failed(response)
        return response.json()

    def get_job_metadata(self, execution_date, project, job) -> dict:
        url = '{optimus_host}/api/v1beta1/project/{project_name}/job/{job_name}/instance'.format(optimus_host=self.host,
                                                                                                 project_name=project,
                                                                                                 job_name=job)
        request_data = {
            "scheduled_at": execution_date,
            "instance_type": "TYPE_TASK",
            "instance_name": "none"
        }
        response = requests.post(url, data=json.dumps(request_data))
        self._raise_error_if_request_failed(response)
        return response.json()

    def notify_event(self, project, namespace, job, event) -> dict:
        url = '{optimus_host}/api/v1beta1/project/{project_name}/namespace/{namespace}/job/{job_name}/event'.format(
            optimus_host=self.host,
            project_name=project,
            namespace=namespace,
            job_name=job,
        )
        request_data = {
            "event": event
        }
        response = requests.post(url, data=json.dumps(request_data))
        self._raise_error_if_request_failed(response)
        return response.json()

    def _raise_error_if_request_failed(self, response):
        if response.status_code != 200:
            log.error("Request to optimus returned non-200 status code. Server response:\n")
            log.error(response.json())
            raise AssertionError("request to optimus returned non-200 status code. url: " + response.url)


class JobSpecTaskWindow:
    def __init__(self, size: str, offset: str, truncate_to: str, optimus_client: OptimusAPIClient):
        self.size = size
        self.offset = offset
        self.truncate_to = truncate_to
        self._optimus_client = optimus_client

    def get(self, scheduled_at: str) -> (datetime, datetime):
        api_response = self._fetch_task_window(scheduled_at)
        return (
            self._parse_datetime(api_response['start']),
            self._parse_datetime(api_response['end']),
        )
    
    # window start is inclusive 
    def get_schedule_window(self, scheduled_at: str , upstream_schedule: str) -> (str, str):
        api_response = self._fetch_task_window(scheduled_at)

        schedule_time_window_start = self._parse_datetime(api_response['start'])
        schedule_time_window_end   = self._parse_datetime(api_response['end'])

        job_cron_iter = croniter(upstream_schedule, schedule_time_window_start)
        schedule_time_window_inclusive_start = job_cron_iter.get_next(datetime)
        return (
            self._parse_datetime_utc_str(schedule_time_window_inclusive_start),
            self._parse_datetime_utc_str(schedule_time_window_end),
        )

    def _parse_datetime(self, timestamp):
        return datetime.strptime(timestamp, TIMESTAMP_FORMAT)

    def _parse_datetime_utc_str(self, timestamp):
        return timestamp.strftime(TIMESTAMP_FORMAT)

    def _fetch_task_window(self, scheduled_at: str) -> dict:
        return self._optimus_client.get_task_window(scheduled_at, self.size, self.offset, self.truncate_to)


class SuperExternalTaskSensor(BaseSensorOperator):

    @apply_defaults
    def __init__(
            self,
            optimus_hostname: str,
            upstream_optimus_project: str,
            upstream_optimus_job: str,
            window_size: str,
            *args,
            **kwargs) -> None:
        kwargs['mode'] = kwargs.get('mode', 'reschedule')
        super().__init__(**kwargs)
        self.optimus_project = upstream_optimus_project
        self.optimus_job = upstream_optimus_job
        self.window_size = window_size
        self._optimus_client = OptimusAPIClient(optimus_hostname)

    def poke(self, context):
        schedule_time = context['next_execution_date']

        # parse relevant metadata from the job metadata to build the task window
        # TODO this needs to be updated to use optimus get job spec
        upstream_schedule = self.get_schedule_interval(schedule_time)

        last_upstream_schedule_time, _ = self.get_last_upstream_times(
            schedule_time, upstream_schedule)

        # get schedule window
        task_window = JobSpecTaskWindow(self.window_size, 0, "m", self._optimus_client)
        schedule_time_window_start, schedule_time_window_end = task_window.get_schedule_window(
            last_upstream_schedule_time.strftime(TIMESTAMP_FORMAT),upstream_schedule)
        

        self.log.info("waiting for upstream runs between: {} - {} schedule times of airflow dag run".format(
            schedule_time_window_start, schedule_time_window_end))

        if not self._are_all_job_runs_successful(schedule_time_window_start, schedule_time_window_end):
            self.log.warning("unable to find enough successful executions for upstream '{}' in "
                             "'{}' dated between {} and {}(inclusive), rescheduling sensor".
                             format(self.optimus_job, self.optimus_project, schedule_time_window_start,
                                    schedule_time_window_end))
            return False
        return True

    def get_last_upstream_times(self, schedule_time_of_current_job, upstream_schedule_interval):
        second_ahead_of_schedule_time = schedule_time_of_current_job + timedelta(seconds=1)
        c = croniter(upstream_schedule_interval, second_ahead_of_schedule_time)
        last_upstream_schedule_time = c.get_prev(datetime)
        last_upstream_execution_date = c.get_prev(datetime)
        return last_upstream_schedule_time, last_upstream_execution_date

    def get_schedule_interval(self, schedule_time):
        schedule_time_str = schedule_time.strftime(TIMESTAMP_FORMAT)
        job_metadata = self._optimus_client.get_job_metadata(schedule_time_str, self.optimus_project, self.optimus_job)
        upstream_schedule = lookup_non_standard_cron_expression(job_metadata['job']['interval'])
        return upstream_schedule

    # TODO the api will be updated with getJobRuns even though the field here refers to scheduledAt
    #  it points to execution_date
    def _are_all_job_runs_successful(self, schedule_time_window_start, schedule_time_window_end) -> bool:
        api_response = self._optimus_client.get_job_run(self.optimus_project, self.optimus_job, schedule_time_window_start, schedule_time_window_end)
        self.log.info("job_run api response :: {}".format(api_response))
        for job_run in api_response['jobRuns']:
            if job_run['state'] != 'success':
                self.log.info("failed for run :: {}".format(job_run))
                return False
        return True

    def _parse_datetime(self, timestamp) -> datetime:
        try:
            return datetime.strptime(timestamp, TIMESTAMP_FORMAT)
        except ValueError:
            return datetime.strptime(timestamp, TIMESTAMP_MS_FORMAT)


def optimus_failure_notify(context):
    params = context.get("params")
    optimus_client = OptimusAPIClient(params["optimus_hostname"])

    taskfail_alert = int(Variable.get("taskfail_alert", default_var=1))
    if taskfail_alert != 1:
        return "suppressed failure alert"

    current_dag_id = context.get('task_instance').dag_id
    current_execution_date = context.get('execution_date')

    # failure message pushed by failed tasks
    failure_messages = []

    def _xcom_value_has_error(_xcom) -> bool:
        return _xcom.key == XCOM_RETURN_KEY and isinstance(_xcom.value, dict) and 'error' in _xcom.value and \
               _xcom.value['error'] is not None

    for xcom in XCom.get_many(
            current_execution_date,
            key=None,
            task_ids=None,
            dag_ids=current_dag_id,
            include_prior_dates=False,
            limit=10):
        if xcom.key == 'error':
            failure_messages.append(xcom.value)
        if _xcom_value_has_error(xcom):
            failure_messages.append(xcom.value['error'])
    failure_message = ", ".join(failure_messages)
    print("failures: {}".format(failure_message))

    message = {
        "log_url": context.get('task_instance').log_url,
        "task_id": context.get('task_instance').task_id,
        "run_id": context.get('run_id'),
        "duration": str(context.get('task_instance').duration),
        "message": failure_message,
        "exception": str(context.get('exception')),
        "scheduled_at": current_execution_date.strftime(TIMESTAMP_FORMAT)
    }
    event = {
        "type": "TYPE_FAILURE",
        "value": message,
    }
    # post event
    resp = optimus_client.notify_event(params["project_name"], params["namespace"], params["job_name"], event)
    print("posted event ", params, event, resp)
    return


def optimus_sla_miss_notify(dag, task_list, blocking_task_list, slas, blocking_tis):
    params = dag.params
    optimus_client = OptimusAPIClient(params["optimus_hostname"])

    slamiss_alert = int(Variable.get("slamiss_alert", default_var=1))
    if slamiss_alert != 1:
        return "suppressed slamiss alert"

    sla_list = []
    for sla in slas:
        sla_list.append({
            'task_id': sla.task_id,
            'dag_id': sla.dag_id,
            'scheduled_at': sla.execution_date.strftime(TIMESTAMP_FORMAT),
            'timestamp': sla.timestamp.strftime(TIMESTAMP_FORMAT)
        })

    current_dag_id = dag.dag_id
    webserver_url = conf.get(section='webserver', key='base_url')
    message = {
        "slas": sla_list,
        "job_url": "{}/tree?dag_id={}".format(webserver_url, current_dag_id),
    }

    event = {
        "type": "TYPE_SLA_MISS",
        "value": message,
    }
    # post event
    resp = optimus_client.notify_event(params["project_name"], params["namespace"], params["job_name"], event)
    print("posted event ", params, event, resp)
    return


# everything below this is here for legacy reasons, should be cleaned up in future

def alert_failed_to_slack(context):
    SLACK_CONN_ID = "slack_alert"
    TASKFAIL_ALERT = int(Variable.get("taskfail_alert", default_var=1))
    SLACK_CHANNEL = Variable.get("slack_channel")

    def _xcom_value_has_error(_xcom) -> bool:
        return _xcom.key == XCOM_RETURN_KEY and isinstance(_xcom.value, dict) and 'error' in _xcom.value and \
               _xcom.value['error'] != None

    if TASKFAIL_ALERT != 1:
        return "suppressed failure alert"

    slack_token = ""
    try:
        slack_token = BaseHook.get_connection(SLACK_CONN_ID).password
    except:
        print("no slack connection variable set")
        return "{connection} connection variable not defined, unable to send alerts".format(connection=SLACK_CONN_ID)

    if not SLACK_CHANNEL:
        return "no slack channel variable set"

    current_dag_id = context.get('task_instance').dag_id
    current_task_id = context.get('task_instance').task_id
    current_execution_date = context.get('execution_date')

    # failure message pushed by failed tasks
    failure_messages = []
    for xcom in XCom.get_many(
            current_execution_date,
            key=None,
            task_ids=None,
            dag_ids=current_dag_id,
            include_prior_dates=False,
            limit=10):
        if xcom.key == 'error':
            failure_messages.append(xcom.value)
        if _xcom_value_has_error(xcom):
            failure_messages.append(xcom.value['error'])
    failure_message = ", ".join(failure_messages)
    print("failures: {}".format(failure_message))

    message_body = "\n".join([
        "• *DAG*: {}".format(current_dag_id),
        "• *Task*: {}".format(current_task_id),
        "• *Execution Time*: {}".format(current_execution_date),
        "• *Run ID*: {}".format(context.get('run_id'))
    ])

    message_footer = "\n".join([
        ":blob-facepalm: Owner: {}".format(context.get('dag').owner),
        ":hourglass: Duration: {} sec".format(context.get('task_instance').duration),
        ":memo: Details: {}".format(failure_message)
    ])

    blocks = [
        {
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": "Task failed :fire:"
            }
        },
        {
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": message_body
            }
        },
        {
            "type": "actions",
            "elements": [
                {
                    "type": "button",
                    "style": "danger",
                    "text": {
                        "type": "plain_text",
                        "text": "View log :airflow:",
                    },
                    "url": context.get('task_instance').log_url,
                    "action_id": "view_log",
                }
            ]
        },
        {
            "type": "divider"
        },
        {
            "type": "context",
            "elements": [
                {
                    "type": "mrkdwn",
                    "text": message_footer
                },
            ]
        },
    ]
    failed_alert = SlackAPIPostOperator(
        slack_conn_id=SLACK_CONN_ID,
        token=slack_token,
        blocks=blocks,
        task_id='slack_failed_alert',
        channel=SLACK_CHANNEL
    )
    return failed_alert.execute(context=context)


class ExternalHttpSensor(BaseSensorOperator):
    """
    Executes a HTTP GET statement and returns False on failure caused by
    404 Not Found

    :param method: The HTTP request method to use
    :param endpoint: The relative part of the full url
    :param request_params: The parameters to be added to the GET url
    :param headers: The HTTP headers to be added to the GET request

    """

    template_fields = ('endpoint', 'request_params', 'headers')

    def __init__(
            self,
            endpoint: str,
            method: str = 'GET',
            request_params: Optional[Dict[str, Any]] = None,
            headers: Optional[Dict[str, Any]] = None,
            *args,
            **kwargs,
    ) -> None:
        kwargs['mode'] = kwargs.get('mode', 'reschedule')
        super().__init__(**kwargs)
        self.endpoint = endpoint
        self.request_params = request_params or {}
        self.headers = headers or {}

    def poke(self, context: 'Context') -> bool:
        self.log.info('Poking: %s', self.endpoint)
        r = requests.get(url=self.endpoint, headers=self.headers, params=self.request_params)
        if (r.status_code >= 200 and r.status_code <= 300):
            return True
        return False
