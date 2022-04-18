# Code generated by optimus dev. DO NOT EDIT.

from typing import Any, Callable, Dict, Optional
from datetime import datetime, timedelta, timezone

from airflow.models import DAG, Variable, DagRun, DagModel, TaskInstance, BaseOperator, XCom, XCOM_RETURN_KEY
from airflow.kubernetes.secret import Secret
from airflow.configuration import conf
from airflow.utils.weight_rule import WeightRule
from kubernetes.client import models as k8s


from __lib import optimus_failure_notify, optimus_sla_miss_notify, SuperKubernetesPodOperator, \
    SuperExternalTaskSensor, ExternalHttpSensor

SENSOR_DEFAULT_POKE_INTERVAL_IN_SECS = int(Variable.get("sensor_poke_interval_in_secs", default_var=15 * 60))
SENSOR_DEFAULT_TIMEOUT_IN_SECS = int(Variable.get("sensor_timeout_in_secs", default_var=15 * 60 * 60))
DAG_RETRIES = int(Variable.get("dag_retries", default_var=3))
DAG_RETRY_DELAY = int(Variable.get("dag_retry_delay_in_secs", default_var=5 * 60))
DAGRUN_TIMEOUT_IN_SECS = int(Variable.get("dagrun_timeout_in_secs", default_var=3 * 24 * 60 * 60))

default_args = {
    "params": {
        "project_name": "foo-project",
        "namespace": "bar-namespace",
        "job_name": "foo",
        "optimus_hostname": "http://airflow.example.io"
    },
    "owner": "mee@mee",
    "depends_on_past": False,
    "retries": 4,
    "retry_delay": timedelta(seconds=DAG_RETRY_DELAY),
    "retry_exponential_backoff": True,
    "priority_weight": 2000,
    "start_date": datetime.strptime("2000-11-11T00:00:00", "%Y-%m-%dT%H:%M:%S"),
    "end_date": datetime.strptime("2020-11-11T00:00:00","%Y-%m-%dT%H:%M:%S"),
    "on_failure_callback": optimus_failure_notify,
    "weight_rule": WeightRule.ABSOLUTE
}

dag = DAG(
    dag_id="foo",
    default_args=default_args,
    schedule_interval="* * * * *",
    sla_miss_callback=optimus_sla_miss_notify,
    catchup=True,
    dagrun_timeout=timedelta(seconds=DAGRUN_TIMEOUT_IN_SECS),
    tags = [
            "optimus",
           ]
)

transformation_secret = Secret(
    "volume",
    "/opt/optimus/secrets",
    "optimus-task-bq",
    "auth.json"
)

transformation_bq = SuperKubernetesPodOperator(
    image_pull_policy="Always",
    namespace = conf.get('kubernetes', 'namespace', fallback="default"),
    image = "example.io/namespace/image:latest",
    cmds=[],
    name="bq",
    task_id="bq",
    get_logs=True,
    dag=dag,
    in_cluster=True,
    is_delete_operator_pod=True,
    do_xcom_push=False,
    secrets=[transformation_secret],
    env_vars = [
        k8s.V1EnvVar(name="JOB_NAME",value='foo'),
        k8s.V1EnvVar(name="OPTIMUS_HOST",value='http://airflow.example.io'),
        k8s.V1EnvVar(name="JOB_LABELS",value='orchestrator=optimus'),
        k8s.V1EnvVar(name="JOB_DIR",value='/data'),
        k8s.V1EnvVar(name="PROJECT",value='foo-project'),
        k8s.V1EnvVar(name="NAMESPACE",value='bar-namespace'),
        k8s.V1EnvVar(name="INSTANCE_TYPE",value='task'),
        k8s.V1EnvVar(name="INSTANCE_NAME",value='bq'),
        k8s.V1EnvVar(name="SCHEDULED_AT",value='{{ next_execution_date }}'),
    ],
    sla=timedelta(seconds=7200),
    reattach_on_restart=True
)

# hooks loop start

hook_transporter_secret = Secret(
    "volume",
    "/opt/optimus/secrets",
    "optimus-hook-transporter",
    "auth.json"
)

hook_transporter = SuperKubernetesPodOperator(
    image_pull_policy="Always",
    namespace = conf.get('kubernetes', 'namespace', fallback="default"),
    image = "example.io/namespace/hook-image:latest",
    cmds=[],
    name="hook_transporter",
    task_id="hook_transporter",
    get_logs=True,
    dag=dag,
    in_cluster=True,
    is_delete_operator_pod=True,
    do_xcom_push=False,
    secrets=[hook_transporter_secret],
    env_vars = [
        k8s.V1EnvVar(name="JOB_NAME",value='foo'),
        k8s.V1EnvVar(name="OPTIMUS_HOST",value='http://airflow.example.io'),
        k8s.V1EnvVar(name="JOB_LABELS",value='orchestrator=optimus'),
        k8s.V1EnvVar(name="JOB_DIR",value='/data'),
        k8s.V1EnvVar(name="PROJECT",value='foo-project'),
        k8s.V1EnvVar(name="NAMESPACE",value='bar-namespace'),
        k8s.V1EnvVar(name="INSTANCE_TYPE",value='hook'),
        k8s.V1EnvVar(name="INSTANCE_NAME",value='transporter'),
        k8s.V1EnvVar(name="SCHEDULED_AT",value='{{ next_execution_date }}'),
        # rest of the env vars are pulled from the container by making a GRPC call to optimus
    ],
    reattach_on_restart=True
)


hook_predator = SuperKubernetesPodOperator(
    image_pull_policy="Always",
    namespace = conf.get('kubernetes', 'namespace', fallback="default"),
    image = "example.io/namespace/predator-image:latest",
    cmds=[],
    name="hook_predator",
    task_id="hook_predator",
    get_logs=True,
    dag=dag,
    in_cluster=True,
    is_delete_operator_pod=True,
    do_xcom_push=False,
    secrets=[],
    env_vars = [
        k8s.V1EnvVar(name="JOB_NAME",value='foo'),
        k8s.V1EnvVar(name="OPTIMUS_HOST",value='http://airflow.example.io'),
        k8s.V1EnvVar(name="JOB_LABELS",value='orchestrator=optimus'),
        k8s.V1EnvVar(name="JOB_DIR",value='/data'),
        k8s.V1EnvVar(name="PROJECT",value='foo-project'),
        k8s.V1EnvVar(name="NAMESPACE",value='bar-namespace'),
        k8s.V1EnvVar(name="INSTANCE_TYPE",value='hook'),
        k8s.V1EnvVar(name="INSTANCE_NAME",value='predator'),
        k8s.V1EnvVar(name="SCHEDULED_AT",value='{{ next_execution_date }}'),
        # rest of the env vars are pulled from the container by making a GRPC call to optimus
    ],
    reattach_on_restart=True
)


hook_hook__dash__for__dash__fail = SuperKubernetesPodOperator(
    image_pull_policy="Always",
    namespace = conf.get('kubernetes', 'namespace', fallback="default"),
    image = "example.io/namespace/fail-image:latest",
    cmds=[],
    name="hook_hook-for-fail",
    task_id="hook_hook-for-fail",
    get_logs=True,
    dag=dag,
    in_cluster=True,
    is_delete_operator_pod=True,
    do_xcom_push=False,
    secrets=[],
    env_vars = [
        k8s.V1EnvVar(name="JOB_NAME",value='foo'),
        k8s.V1EnvVar(name="OPTIMUS_HOST",value='http://airflow.example.io'),
        k8s.V1EnvVar(name="JOB_LABELS",value='orchestrator=optimus'),
        k8s.V1EnvVar(name="JOB_DIR",value='/data'),
        k8s.V1EnvVar(name="PROJECT",value='foo-project'),
        k8s.V1EnvVar(name="NAMESPACE",value='bar-namespace'),
        k8s.V1EnvVar(name="INSTANCE_TYPE",value='hook'),
        k8s.V1EnvVar(name="INSTANCE_NAME",value='hook-for-fail'),
        k8s.V1EnvVar(name="SCHEDULED_AT",value='{{ next_execution_date }}'),
        # rest of the env vars are pulled from the container by making a GRPC call to optimus
    ],
    trigger_rule="one_failed",
    reattach_on_restart=True
)
# hooks loop ends


# create upstream sensors

wait_foo__dash__intra__dash__dep__dash__job = SuperExternalTaskSensor(
    optimus_hostname="http://airflow.example.io",
    upstream_optimus_project="foo-project",
    upstream_optimus_job="foo-intra-dep-job",
    window_size="1h0m0s",
    poke_interval=SENSOR_DEFAULT_POKE_INTERVAL_IN_SECS,
    timeout=SENSOR_DEFAULT_TIMEOUT_IN_SECS,
    task_id="wait_foo-intra-dep-job-bq",
    dag=dag
)
wait_foo__dash__inter__dash__dep__dash__job = SuperExternalTaskSensor(
    optimus_hostname="http://airflow.example.io",
    upstream_optimus_project="foo-external-project",
    upstream_optimus_job="foo-inter-dep-job",
    window_size="1h0m0s",
    poke_interval=SENSOR_DEFAULT_POKE_INTERVAL_IN_SECS,
    timeout=SENSOR_DEFAULT_TIMEOUT_IN_SECS,
    task_id="wait_foo-inter-dep-job-bq",
    dag=dag
)

# arrange inter task dependencies
####################################

# upstream sensors -> base transformation task
wait_foo__dash__intra__dash__dep__dash__job >> transformation_bq
wait_foo__dash__inter__dash__dep__dash__job >> transformation_bq

# set inter-dependencies between task and hooks
hook_transporter >> transformation_bq
transformation_bq >> hook_predator
transformation_bq >> hook_hook__dash__for__dash__fail

# set inter-dependencies between hooks and hooks
hook_transporter >> hook_predator

# arrange failure hook after post hooks

hook_predator >> [ hook_hook__dash__for__dash__fail,]