package v1beta1_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/odpf/salt/log"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	v1 "github.com/odpf/optimus/api/handler/v1beta1"
	pb "github.com/odpf/optimus/api/proto/odpf/optimus/core/v1beta1"
	"github.com/odpf/optimus/core/progress"
	"github.com/odpf/optimus/mock"
	"github.com/odpf/optimus/models"
)

type JobSpecServiceServerTestSuite struct {
	suite.Suite
	ctx              context.Context //nolint:containedctx
	projectService   *mock.ProjectService
	namespaceService *mock.NamespaceService
	jobService       *mock.JobService // TODO: refactor to service package
	adapter          *mock.ProtoAdapter
	log              log.Logger
	progressObserver progress.Observer

	jobReq        *pb.DeployJobSpecificationRequest
	resourceReq   *pb.DeployResourceSpecificationRequest
	projectSpec   models.ProjectSpec
	namespaceSpec models.NamespaceSpec
}

func (s *JobSpecServiceServerTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.namespaceService = new(mock.NamespaceService)
	s.adapter = new(mock.ProtoAdapter)
	s.jobService = new(mock.JobService)
	s.log = log.NewNoop()

	s.projectSpec = models.ProjectSpec{}
	s.projectSpec.Name = "project-a"
	s.projectSpec.ID = models.ProjectID(uuid.MustParse("26a0d6a0-13c6-4b30-ae6f-29233df70f31"))

	s.namespaceSpec = models.NamespaceSpec{}
	s.namespaceSpec.Name = "ns1"
	s.namespaceSpec.ID = uuid.MustParse("ceba7919-e07d-48b4-a4ce-141d79a3b59d")

	s.jobReq = &pb.DeployJobSpecificationRequest{}
	s.jobReq.ProjectName = s.projectSpec.Name
	s.jobReq.NamespaceName = s.namespaceSpec.Name

	s.resourceReq = &pb.DeployResourceSpecificationRequest{}
	s.resourceReq.DatastoreName = "datastore-1"
	s.resourceReq.ProjectName = s.projectSpec.Name
	s.resourceReq.NamespaceName = s.namespaceSpec.Name
}

func (s *JobSpecServiceServerTestSuite) newJobSpecServiceServer() *v1.JobSpecServiceServer {
	return v1.NewJobSpecServiceServer(
		s.log,
		s.jobService,
		s.adapter,
		s.projectService,
		s.namespaceService,
		s.progressObserver,
	)
}

func TestRuntimeServiceServerJobTestSuite(t *testing.T) {
	s := new(JobSpecServiceServerTestSuite)
	suite.Run(t, s)
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Success_NoJobSpec() {
	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Context").Return(s.ctx)
	stream.On("Recv").Return(s.jobReq, nil).Once()
	stream.On("Recv").Return(nil, io.EOF).Once()

	s.namespaceService.On("Get", s.ctx, s.jobReq.GetProjectName(), s.jobReq.GetNamespaceName()).Return(s.namespaceSpec, nil).Once()
	s.jobService.On("KeepOnly", s.ctx, s.namespaceSpec, mock2.Anything, mock2.Anything).Return(nil).Once()
	s.jobService.On("Sync", s.ctx, s.namespaceSpec, mock2.Anything).Return(nil).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().NoError(err)
	stream.AssertExpectations(s.T())
	s.namespaceService.AssertExpectations(s.T())
	s.jobService.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Success_TwoJobSpecs() {
	jobSpecs := []*pb.JobSpecification{}
	jobSpecs = append(jobSpecs, &pb.JobSpecification{Name: "job-1"})
	jobSpecs = append(jobSpecs, &pb.JobSpecification{Name: "job-2"})
	s.jobReq.Jobs = jobSpecs
	adaptedJobs := []models.JobSpec{}
	adaptedJobs = append(adaptedJobs, models.JobSpec{Name: "job-1"})
	adaptedJobs = append(adaptedJobs, models.JobSpec{Name: "job-2"})

	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Context").Return(s.ctx)
	stream.On("Recv").Return(s.jobReq, nil).Once()
	stream.On("Recv").Return(nil, io.EOF).Once()

	s.namespaceService.On("Get", s.ctx, s.jobReq.GetProjectName(), s.jobReq.GetNamespaceName()).Return(s.namespaceSpec, nil).Once()
	for i := range jobSpecs {
		s.adapter.On("FromJobProto", jobSpecs[i]).Return(adaptedJobs[i], nil).Once()
		s.jobService.On("Create", s.ctx, s.namespaceSpec, adaptedJobs[i]).Return(nil).Once()
	}
	s.jobService.On("KeepOnly", s.ctx, s.namespaceSpec, adaptedJobs, mock2.Anything).Return(nil).Once()
	s.jobService.On("Sync", s.ctx, s.namespaceSpec, mock2.Anything).Return(nil).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().NoError(err)
	stream.AssertExpectations(s.T())
	s.adapter.AssertExpectations(s.T())
	s.namespaceService.AssertExpectations(s.T())
	s.jobService.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Fail_StreamRecvError() {
	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Recv").Return(nil, errors.New("any error")).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().Error(err)
	stream.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Fail_NamespaceServiceGetError() {
	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Context").Return(s.ctx)
	stream.On("Recv").Return(s.jobReq, nil).Once()
	stream.On("Recv").Return(nil, io.EOF).Once()

	s.namespaceService.On("Get", s.ctx, s.jobReq.GetProjectName(), s.jobReq.GetNamespaceName()).Return(models.NamespaceSpec{}, errors.New("any error")).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().Error(err)
	stream.AssertExpectations(s.T())
	s.namespaceService.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Fail_AdapterFromJobProtoError() {
	jobSpecs := []*pb.JobSpecification{}
	jobSpecs = append(jobSpecs, &pb.JobSpecification{Name: "job-1"})
	s.jobReq.Jobs = jobSpecs

	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Context").Return(s.ctx)
	stream.On("Recv").Return(s.jobReq, nil).Once()
	stream.On("Recv").Return(nil, io.EOF).Once()

	s.namespaceService.On("Get", s.ctx, s.jobReq.GetProjectName(), s.jobReq.GetNamespaceName()).Return(s.namespaceSpec, nil).Once()
	s.adapter.On("FromJobProto", jobSpecs[0]).Return(models.JobSpec{}, errors.New("any error")).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().Error(err)
	stream.AssertExpectations(s.T())
	s.adapter.AssertExpectations(s.T())
	s.namespaceService.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Fail_JobServiceCreateError() {
	jobSpecs := []*pb.JobSpecification{}
	jobSpecs = append(jobSpecs, &pb.JobSpecification{Name: "job-1"})
	s.jobReq.Jobs = jobSpecs
	adaptedJobs := []models.JobSpec{}
	adaptedJobs = append(adaptedJobs, models.JobSpec{Name: "job-1"})

	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Context").Return(s.ctx)
	stream.On("Recv").Return(s.jobReq, nil).Once()
	stream.On("Recv").Return(nil, io.EOF).Once()

	s.namespaceService.On("Get", s.ctx, s.jobReq.GetProjectName(), s.jobReq.GetNamespaceName()).Return(s.namespaceSpec, nil).Once()
	s.adapter.On("FromJobProto", jobSpecs[0]).Return(adaptedJobs[0], nil).Once()
	s.jobService.On("Create", s.ctx, s.namespaceSpec, adaptedJobs[0]).Return(errors.New("any error")).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().Error(err)
	stream.AssertExpectations(s.T())
	s.adapter.AssertExpectations(s.T())
	s.namespaceService.AssertExpectations(s.T())
	s.jobService.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Fail_JobServiceKeepOnlyError() {
	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Context").Return(s.ctx)
	stream.On("Recv").Return(s.jobReq, nil).Once()
	stream.On("Recv").Return(nil, io.EOF).Once()

	s.namespaceService.On("Get", s.ctx, s.jobReq.GetProjectName(), s.jobReq.GetNamespaceName()).Return(s.namespaceSpec, nil).Once()
	s.jobService.On("KeepOnly", s.ctx, s.namespaceSpec, mock2.Anything, mock2.Anything).Return(errors.New("any error")).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().Error(err)
	stream.AssertExpectations(s.T())
	s.namespaceService.AssertExpectations(s.T())
	s.jobService.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestDeployJobSpecification_Fail_JobServiceSyncError() {
	stream := new(mock.DeployJobSpecificationServer)
	stream.On("Context").Return(s.ctx)
	stream.On("Recv").Return(s.jobReq, nil).Once()
	stream.On("Recv").Return(nil, io.EOF).Once()

	s.namespaceService.On("Get", s.ctx, s.jobReq.GetProjectName(), s.jobReq.GetNamespaceName()).Return(s.namespaceSpec, nil).Once()
	s.jobService.On("KeepOnly", s.ctx, s.namespaceSpec, mock2.Anything, mock2.Anything).Return(nil).Once()
	s.jobService.On("Sync", s.ctx, s.namespaceSpec, mock2.Anything).Return(errors.New("any error")).Once()
	stream.On("Send", mock2.Anything).Return(nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	err := runtimeServiceServer.DeployJobSpecification(stream)

	s.Assert().Error(err)
	stream.AssertExpectations(s.T())
	s.namespaceService.AssertExpectations(s.T())
	s.jobService.AssertExpectations(s.T())
}

func (s *JobSpecServiceServerTestSuite) TestGetJobSpecification_Success() {
	req := &pb.GetJobSpecificationRequest{}
	req.ProjectName = s.projectSpec.Name
	req.NamespaceName = s.namespaceSpec.Name
	req.JobName = "job-1"
	jobSpec := models.JobSpec{Name: req.JobName}
	jobSpecProto := &pb.JobSpecification{Name: req.JobName}

	s.namespaceService.On("Get", s.ctx, req.ProjectName, req.NamespaceName).Return(s.namespaceSpec, nil).Once()
	s.jobService.On("GetByName", s.ctx, req.JobName, s.namespaceSpec).Return(jobSpec, nil).Once()
	s.adapter.On("ToJobProto", jobSpec).Return(jobSpecProto, nil).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	resp, err := runtimeServiceServer.GetJobSpecification(s.ctx, req)

	s.Assert().NoError(err)
	s.Assert().NotNil(resp)
}

func (s *JobSpecServiceServerTestSuite) TestGetJobSpecification_Fail_NamespaceServiceGetError() {
	req := &pb.GetJobSpecificationRequest{}
	req.ProjectName = s.projectSpec.Name
	req.NamespaceName = s.namespaceSpec.Name
	req.JobName = "job-1"

	s.namespaceService.On("Get", s.ctx, req.ProjectName, req.NamespaceName).Return(models.NamespaceSpec{}, errors.New("any error")).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	resp, err := runtimeServiceServer.GetJobSpecification(s.ctx, req)

	s.Assert().Error(err)
	s.Assert().Nil(resp)
}

func (s *JobSpecServiceServerTestSuite) TestGetJobSpecification_Fail_JobServiceGetByNameError() {
	req := &pb.GetJobSpecificationRequest{}
	req.ProjectName = s.projectSpec.Name
	req.NamespaceName = s.namespaceSpec.Name
	req.JobName = "job-1"

	s.namespaceService.On("Get", s.ctx, req.ProjectName, req.NamespaceName).Return(s.namespaceSpec, nil).Once()
	s.jobService.On("GetByName", s.ctx, req.JobName, s.namespaceSpec).Return(models.JobSpec{}, errors.New("any error")).Once()

	runtimeServiceServer := s.newJobSpecServiceServer()
	resp, err := runtimeServiceServer.GetJobSpecification(s.ctx, req)

	s.Assert().Error(err)
	s.Assert().Nil(resp)
}

// TODO: refactor to test suite
func TestJobSpecificationOnServer(t *testing.T) {
	log := log.NewNoop()
	ctx := context.Background()
	t.Run("RegisterJobSpecification", func(t *testing.T) {
		t.Run("should save a job specification", func(t *testing.T) {
			projectName := "a-data-project"

			projectSpec := models.ProjectSpec{
				Name: projectName,
				Config: map[string]string{
					"BUCKET": "gs://some_folder",
				},
			}

			namespaceSpec := models.NamespaceSpec{
				Name:        "dev-test-namespace-1",
				ProjectSpec: projectSpec,
			}

			jobName := "my-job"
			taskName := "bq2bq"
			execUnit1 := new(mock.BasePlugin)
			execUnit1.On("PluginInfo").Return(&models.PluginInfoResponse{
				Name:  taskName,
				Image: "random-image",
			}, nil)
			defer execUnit1.AssertExpectations(t)

			pluginRepo := new(mock.SupportedPluginRepo)
			pluginRepo.On("GetByName", taskName).Return(&models.Plugin{
				Base: execUnit1,
			}, nil)
			adapter := v1.NewAdapter(pluginRepo, nil)

			jobSpec := models.JobSpec{
				Name: jobName,
				Task: models.JobSpecTask{
					Unit: &models.Plugin{
						Base: execUnit1,
					},
					Config: models.JobSpecConfigs{
						{
							Name:  "DO",
							Value: "THIS",
						},
					},
					Window: models.JobSpecTaskWindow{
						Size:       time.Hour,
						Offset:     0,
						TruncateTo: "d",
					},
				},
				Assets: *models.JobAssets{}.New(
					[]models.JobSpecAsset{
						{
							Name:  "query.sql",
							Value: "select * from 1",
						},
					}),
				Dependencies: map[string]models.JobSpecDependency{},
			}

			jobSvc := new(mock.JobService)
			jobSvc.On("Create", ctx, namespaceSpec, jobSpec).Return(nil)
			jobSvc.On("Check", ctx, namespaceSpec, []models.JobSpec{jobSpec}, mock2.Anything).Return(nil)
			jobSvc.On("Sync", mock2.Anything, namespaceSpec, mock2.Anything).Return(nil)
			defer jobSvc.AssertExpectations(t)

			namespaceService := new(mock.NamespaceService)
			namespaceService.On("Get", ctx, projectSpec.Name, namespaceSpec.Name).Return(namespaceSpec, nil)
			defer namespaceService.AssertExpectations(t)

			jobSpecServiceServer := v1.NewJobSpecServiceServer(
				log,
				jobSvc,
				adapter,
				nil,
				namespaceService,
				nil,
			)

			jobProto := adapter.ToJobProto(jobSpec)
			request := pb.CreateJobSpecificationRequest{
				ProjectName:   projectName,
				NamespaceName: namespaceSpec.Name,
				Spec:          jobProto,
			}
			resp, err := jobSpecServiceServer.CreateJobSpecification(ctx, &request)
			assert.Nil(t, err)
			assert.Equal(t, &pb.CreateJobSpecificationResponse{
				Success: true,
				Message: "job my-job is created and deployed successfully on project a-data-project",
			}, resp)
		})
	})
	t.Run("DeleteJobSpecification", func(t *testing.T) {
		t.Run("should delete the job", func(t *testing.T) {
			projectName := "a-data-project"
			jobName1 := "a-data-job"
			taskName := "a-data-task"

			projectSpec := models.ProjectSpec{
				ID:   models.ProjectID(uuid.New()),
				Name: projectName,
				Config: map[string]string{
					"bucket": "gs://some_folder",
				},
			}

			namespaceSpec := models.NamespaceSpec{
				ID:   uuid.Must(uuid.NewRandom()),
				Name: "dev-test-namespace-1",
				Config: map[string]string{
					"bucket": "gs://some_folder",
				},
				ProjectSpec: projectSpec,
			}

			execUnit1 := new(mock.BasePlugin)
			defer execUnit1.AssertExpectations(t)

			jobSpecs := []models.JobSpec{
				{
					Name: jobName1,
					Task: models.JobSpecTask{
						Unit: &models.Plugin{
							Base: execUnit1,
						},
						Config: models.JobSpecConfigs{
							{
								Name:  "do",
								Value: "this",
							},
						},
					},
					Assets: *models.JobAssets{}.New(
						[]models.JobSpecAsset{
							{
								Name:  "query.sql",
								Value: "select * from 1",
							},
						}),
				},
			}

			pluginRepo := new(mock.SupportedPluginRepo)
			pluginRepo.On("GetByName", taskName).Return(&models.Plugin{
				Base: execUnit1,
			}, nil)
			adapter := v1.NewAdapter(pluginRepo, nil)

			namespaceService := new(mock.NamespaceService)
			namespaceService.On("Get", ctx, projectSpec.Name, namespaceSpec.Name).
				Return(namespaceSpec, nil)
			defer namespaceService.AssertExpectations(t)

			jobSpec := jobSpecs[0]

			jobService := new(mock.JobService)
			jobService.On("GetByName", ctx, jobSpecs[0].Name, namespaceSpec).Return(jobSpecs[0], nil)
			jobService.On("Delete", mock2.Anything, namespaceSpec, jobSpec).Return(nil)
			defer jobService.AssertExpectations(t)

			jobSpecServiceServer := v1.NewJobSpecServiceServer(
				log,
				jobService,
				adapter,
				nil,
				namespaceService,
				nil,
			)

			deployRequest := pb.DeleteJobSpecificationRequest{ProjectName: projectName, JobName: jobSpec.Name, NamespaceName: namespaceSpec.Name}
			resp, err := jobSpecServiceServer.DeleteJobSpecification(ctx, &deployRequest)
			assert.Nil(t, err)
			assert.Equal(t, "job a-data-job has been deleted", resp.GetMessage())
		})
	})

	t.Run("RefreshJobs", func(t *testing.T) {
		t.Run("should refresh jobs successfully", func(t *testing.T) {
			projectName := "a-data-project"
			projectSpec := models.ProjectSpec{
				ID:   models.ProjectID(uuid.New()),
				Name: projectName,
				Config: map[string]string{
					"bucket": "gs://some_folder",
				},
			}
			namespaceSpec := models.NamespaceSpec{
				ID:   uuid.Must(uuid.NewRandom()),
				Name: "dev-test-namespace-1",
				Config: map[string]string{
					"bucket": "gs://some_folder",
				},
				ProjectSpec: projectSpec,
			}
			namespaceNames := []string{namespaceSpec.Name}

			jobSpecRepository := new(mock.JobSpecRepository)
			defer jobSpecRepository.AssertExpectations(t)

			jobSpecRepoFactory := new(mock.JobSpecRepoFactory)
			defer jobSpecRepoFactory.AssertExpectations(t)

			pluginRepo := new(mock.SupportedPluginRepo)
			adapter := v1.NewAdapter(pluginRepo, nil)

			nsService := new(mock.NamespaceService)
			defer nsService.AssertExpectations(t)

			projectJobSpecRepository := new(mock.ProjectJobSpecRepository)
			defer projectJobSpecRepository.AssertExpectations(t)

			projectJobSpecRepoFactory := new(mock.ProjectJobSpecRepoFactory)
			defer projectJobSpecRepoFactory.AssertExpectations(t)

			projectService := new(mock.ProjectService)
			defer projectService.AssertExpectations(t)

			jobService := new(mock.JobService)
			defer jobService.AssertExpectations(t)

			grpcRespStream := new(mock.RefreshJobsServer)
			defer grpcRespStream.AssertExpectations(t)

			jobService.On("Refresh", mock2.Anything, projectSpec.Name, namespaceNames, []string(nil), mock2.Anything).Return(nil)
			grpcRespStream.On("Context").Return(context.Background())

			jobSpecServiceServer := v1.NewJobSpecServiceServer(
				log,
				jobService,
				adapter,
				projectService,
				nsService,
				nil,
			)
			refreshRequest := pb.RefreshJobsRequest{ProjectName: projectName, NamespaceNames: namespaceNames}
			err := jobSpecServiceServer.RefreshJobs(&refreshRequest, grpcRespStream)
			assert.Nil(t, err)
		})
		t.Run("should failed when unable to do refresh jobs", func(t *testing.T) {
			projectName := "a-data-project"
			projectSpec := models.ProjectSpec{
				ID:   models.ProjectID(uuid.New()),
				Name: projectName,
				Config: map[string]string{
					"bucket": "gs://some_folder",
				},
			}
			namespaceSpec := models.NamespaceSpec{
				ID:   uuid.Must(uuid.NewRandom()),
				Name: "dev-test-namespace-1",
				Config: map[string]string{
					"bucket": "gs://some_folder",
				},
				ProjectSpec: projectSpec,
			}
			namespaceNames := []string{namespaceSpec.Name}

			jobSpecRepository := new(mock.JobSpecRepository)
			defer jobSpecRepository.AssertExpectations(t)

			jobSpecRepoFactory := new(mock.JobSpecRepoFactory)
			defer jobSpecRepoFactory.AssertExpectations(t)

			pluginRepo := new(mock.SupportedPluginRepo)
			adapter := v1.NewAdapter(pluginRepo, nil)

			nsService := new(mock.NamespaceService)
			defer nsService.AssertExpectations(t)

			projectJobSpecRepository := new(mock.ProjectJobSpecRepository)
			defer projectJobSpecRepository.AssertExpectations(t)

			projectJobSpecRepoFactory := new(mock.ProjectJobSpecRepoFactory)
			defer projectJobSpecRepoFactory.AssertExpectations(t)

			projectService := new(mock.ProjectService)
			defer projectService.AssertExpectations(t)

			jobService := new(mock.JobService)
			defer jobService.AssertExpectations(t)

			grpcRespStream := new(mock.RefreshJobsServer)
			defer grpcRespStream.AssertExpectations(t)

			errorMsg := "internal error"
			jobService.On("Refresh", mock2.Anything, projectSpec.Name, namespaceNames, []string(nil), mock2.Anything).Return(errors.New(errorMsg))
			grpcRespStream.On("Context").Return(context.Background())

			jobSpecServiceServer := v1.NewJobSpecServiceServer(
				log,
				jobService,
				adapter,
				projectService,
				nsService,
				nil,
			)
			refreshRequest := pb.RefreshJobsRequest{ProjectName: projectName, NamespaceNames: namespaceNames}
			err := jobSpecServiceServer.RefreshJobs(&refreshRequest, grpcRespStream)
			assert.Contains(t, err.Error(), errorMsg)
		})
	})
}
