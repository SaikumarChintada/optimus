package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/odpf/salt/log"
	cli "github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/odpf/optimus/api/proto/odpf/optimus/core/v1beta1"
	"github.com/odpf/optimus/config"
)

const (
	jobStatusTimeout = time.Second * 30
)

func jobRunListCommand(conf *config.ClientConfig) *cli.Command {
	cmd := &cli.Command{
		Use:     "list-runs",
		Short:   "Get Job run details",
		Example: `optimus job runs <sample_job_goes_here> [--project \"project-id\"] [--start_date \"2006-01-02T15:04:05Z07:00\" --end_date \"2006-01-02T15:04:05Z07:00\"]`,
		Args:    cli.MinimumNArgs(1),
	}
	var startDate string
	var endDate string
	cmd.Flags().StringP("project-name", "p", defaultProjectName, "Project name of optimus managed repository")
	cmd.Flags().String("host", defaultHost, "Optimus service endpoint url")
	cmd.Flags().StringVar(&startDate, "start_date", "", "start date of job run")
	cmd.Flags().StringVar(&endDate, "end_date", "", "end date of job run")
	cmd.RunE = func(c *cli.Command, args []string) error {
		projectName := conf.Project.Name
		host := conf.Host
		l := initClientLogger(conf.Log)
		jobName := args[0]
		l.Info(fmt.Sprintf("Requesting status for project %s, job %s from %s",
			projectName, jobName, host))
		return getJobRunList(l, host, projectName, jobName, startDate, endDate)
	}
	return cmd
}

func getJobRunList(l log.Logger, host, projectName, jobName, startDate, endDate string) error {
	var err error
	var req *pb.JobRunRequest
	// validate date if it has given
	err = validateDateArgs(startDate, endDate)
	if err != nil {
		return err
	}
	// create job run grpc request
	req, err = createJobRunRequest(projectName, jobName, startDate, endDate)
	if err != nil {
		return err
	}
	err = callJobRun(l, host, req)
	return err
}

func callJobRun(l log.Logger, host string, jobRunRequest *pb.JobRunRequest) error {
	ctx, conn, closeConn, err := initClientConnection(host, jobStatusTimeout)
	if err != nil {
		return err
	}
	defer closeConn()

	run := pb.NewJobRunServiceClient(conn)
	spinner := NewProgressBar()
	spinner.Start("please wait...")
	jobRunResponse, err := run.JobRun(ctx, jobRunRequest)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("request failed for job %s: %w", jobRunRequest.JobName, err)
	}

	jobRuns := jobRunResponse.GetJobRuns()
	for _, jobRun := range jobRuns {
		l.Info(fmt.Sprintf("%s - %s", jobRun.GetScheduledAt().AsTime(), jobRun.GetState()))
	}
	l.Info(coloredSuccess("\nFound %d jobRun instances.", len(jobRuns)))
	return err
}

func createJobRunRequest(projectName, jobName, startDate, endDate string) (*pb.JobRunRequest, error) {
	var req *pb.JobRunRequest
	if startDate == "" && endDate == "" {
		req = &pb.JobRunRequest{
			ProjectName: projectName,
			JobName:     jobName,
		}
		return req, nil
	}
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return req, fmt.Errorf("start_date %w", err)
	}
	sDate := timestamppb.New(start)
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return req, fmt.Errorf("end_date %w", err)
	}
	eDate := timestamppb.New(end)
	req = &pb.JobRunRequest{
		ProjectName: projectName,
		JobName:     jobName,
		StartDate:   sDate,
		EndDate:     eDate,
	}
	return req, nil
}

func validateDateArgs(startDate, endDate string) error {
	if startDate == "" && endDate != "" {
		return errors.New("please provide the start date")
	}
	if startDate != "" && endDate == "" {
		return errors.New("please provide the end date")
	}
	return nil
}
