package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/odpf/salt/log"
	"github.com/spf13/afero"
	cli "github.com/spf13/cobra"

	v1handler "github.com/odpf/optimus/api/handler/v1beta1"
	pb "github.com/odpf/optimus/api/proto/odpf/optimus/core/v1beta1"
	"github.com/odpf/optimus/config"
	"github.com/odpf/optimus/models"
	"github.com/odpf/optimus/store/local"
)

const (
	runJobTimeout = time.Minute * 1
)

func jobRunCommand(conf *config.ClientConfig) *cli.Command {
	var (
		namespaceName string
		projectName   string
		host          string
	)

	cmd := &cli.Command{
		Use:     "run",
		Short:   "[EXPERIMENTAL] run the provided job on optimus cluster",
		Args:    cli.MinimumNArgs(1),
		Example: "optimus job run <job_name>",
		Hidden:  true,
		RunE: func(c *cli.Command, args []string) error {
			projectName = conf.Project.Name
			host = conf.Host
			l := initClientLogger(conf.Log)
			pluginRepo := models.PluginRegistry
			namespace, err := conf.GetNamespaceByName(namespaceName)
			if err != nil {
				return err
			}
			jobSpecFs := afero.NewBasePathFs(afero.NewOsFs(), namespace.Job.Path)
			jobSpecRepo := local.NewJobSpecRepository(
				jobSpecFs,
				local.NewJobSpecAdapter(pluginRepo),
			)
			jobSpec, err := jobSpecRepo.GetByName(args[0])
			if err != nil {
				return err
			}
			return runJobSpecificationRequest(l, projectName, namespace.Name, host, jobSpec, pluginRepo)
		},
	}
	cmd.Flags().StringVarP(&namespaceName, "namespace", "n", namespaceName, "Namespace of the resource within project")
	cmd.MarkFlagRequired("namespace")
	return cmd
}

func runJobSpecificationRequest(l log.Logger, projectName, namespace, host string, jobSpec models.JobSpec, pluginRepo models.PluginRepository) (err error) {
	ctx, conn, closeConn, err := initClientConnection(host, runJobTimeout)
	if err != nil {
		return err
	}
	defer closeConn()

	adapt := v1handler.NewAdapter(pluginRepo, nil)
	adaptedSpec := adapt.ToJobProto(jobSpec)

	l.Info("please wait...")
	jobRun := pb.NewJobRunServiceClient(conn)
	jobResponse, err := jobRun.RunJob(ctx, &pb.RunJobRequest{
		ProjectName:   projectName,
		NamespaceName: namespace,
		Specifications: []*pb.JobSpecification{
			adaptedSpec,
		},
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			l.Info("process took too long, timing out")
		}
		return fmt.Errorf("request failed for job %s: %w", jobSpec.Name, err)
	}
	l.Info(jobResponse.String())
	return nil
}
