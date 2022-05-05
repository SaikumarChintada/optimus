package exd_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	tMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/odpf/optimus/ext/exd"
	"github.com/odpf/optimus/mock"
)

type ManagerTestSuite struct {
	suite.Suite
}

func (m *ManagerTestSuite) TestInstall() {
	defaultParser := exd.ParseRegistry
	defer func() { exd.ParseRegistry = defaultParser }()
	defaultNewClient := exd.NewClientRegistry
	defer func() { exd.NewClientRegistry = defaultNewClient }()

	m.Run("should return error if remote path is empty", func() {
		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		manifester := &mock.Manifester{}
		installer := &mock.Installer{}

		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		var remotePath string
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if validation is error", func() {
		manager := &exd.Manager{}
		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if error loading manifest", func() {
		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		installer := &mock.Installer{}

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(nil, errors.New("random error"))

		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if error extracting metadata", func() {
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return nil, errors.New("extraction failed")
			},
		}

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		installer := &mock.Installer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if no parser could recognize remote path", func() {
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return nil, exd.ErrUnrecognizedRemotePath
			},
		}

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		installer := &mock.Installer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if error getting new client", func() {
		provider := "testing"
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return &exd.Metadata{
					ProviderName: provider,
				}, nil
			},
		}

		newClientFactory := &exd.NewClientFactory{}
		exd.NewClientRegistry = newClientFactory

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		installer := &mock.Installer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if error getting release", func() {
		provider := "testing"
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return &exd.Metadata{
					ProviderName: provider,
				}, nil
			},
		}

		client := &mock.Client{}
		client.On("GetRelease", tMock.Anything).Return(nil, errors.New("random error"))
		newClientFactory := &exd.NewClientFactory{}
		newClientFactory.Add(provider, func(ctx context.Context, httpDoer exd.HTTPDoer) (exd.Client, error) {
			return client, nil
		})
		exd.NewClientRegistry = newClientFactory

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		installer := &mock.Installer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if remote path is already installed", func() {
		provider := "testing"
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return &exd.Metadata{
					ProviderName: provider,
					OwnerName:    "gojek",
					RepoName:     "optimus-extension-valor",
					TagName:      "",
				}, nil
			},
		}

		client := &mock.Client{}
		client.On("GetRelease", tMock.Anything).Return(&exd.RepositoryRelease{
			Name: "v1.0",
		}, nil)
		newClientFactory := &exd.NewClientFactory{}
		newClientFactory.Add(provider, func(ctx context.Context, httpDoer exd.HTTPDoer) (exd.Client, error) {
			return client, nil
		})
		exd.NewClientRegistry = newClientFactory

		commandName := "valor"
		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{
			RepositoryOwners: []*exd.RepositoryOwner{
				{
					Name:     "gojek",
					Provider: provider,
					Projects: []*exd.RepositoryProject{
						{
							Name:          "optimus-extension-valor",
							CommandName:   commandName,
							ActiveTagName: "v1.0",
						},
					},
				},
			},
		}, nil)

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		installer := &mock.Installer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if error when downloading", func() {
		provider := "testing"
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return &exd.Metadata{
					ProviderName: provider,
				}, nil
			},
		}

		client := &mock.Client{}
		client.On("GetRelease", tMock.Anything).Return(&exd.RepositoryRelease{
			Name: "v1.0",
		}, nil)
		client.On("DownloadAsset", tMock.Anything).Return(nil, errors.New("random error"))
		newClientFactory := &exd.NewClientFactory{}
		newClientFactory.Add(provider, func(ctx context.Context, httpDoer exd.HTTPDoer) (exd.Client, error) {
			return client, nil
		})
		exd.NewClientRegistry = newClientFactory

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		installer := &mock.Installer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if error when preparing installation", func() {
		provider := "testing"
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return &exd.Metadata{
					ProviderName: provider,
				}, nil
			},
		}

		client := &mock.Client{}
		client.On("GetRelease", tMock.Anything).Return(&exd.RepositoryRelease{
			Name: "v1.0",
		}, nil)
		client.On("DownloadAsset", tMock.Anything).Return([]byte{}, nil)
		newClientFactory := &exd.NewClientFactory{}
		newClientFactory.Add(provider, func(ctx context.Context, httpDoer exd.HTTPDoer) (exd.Client, error) {
			return client, nil
		})
		exd.NewClientRegistry = newClientFactory

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)

		installer := &mock.Installer{}
		installer.On("Prepare", tMock.Anything).Return(errors.New("random error"))

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should return error if error when executing installation", func() {
		provider := "testing"
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return &exd.Metadata{
					ProviderName: provider,
				}, nil
			},
		}

		client := &mock.Client{}
		client.On("GetRelease", tMock.Anything).Return(&exd.RepositoryRelease{
			Name: "v1.0",
		}, nil)
		client.On("DownloadAsset", tMock.Anything).Return([]byte{}, nil)
		newClientFactory := &exd.NewClientFactory{}
		newClientFactory.Add(provider, func(ctx context.Context, httpDoer exd.HTTPDoer) (exd.Client, error) {
			return client, nil
		})
		exd.NewClientRegistry = newClientFactory

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)

		installer := &mock.Installer{}
		installer.On("Prepare", tMock.Anything).Return(nil)
		installer.On("Install", tMock.Anything, tMock.Anything).Return(errors.New("random error"))

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.Error(actualErr)
	})

	m.Run("should update manifest and return nil if no error is encountered", func() {
		provider := "testing"
		exd.ParseRegistry = []exd.Parser{
			func(remotePath string) (*exd.Metadata, error) {
				return &exd.Metadata{
					ProviderName: provider,
				}, nil
			},
		}

		client := &mock.Client{}
		client.On("GetRelease", tMock.Anything).Return(&exd.RepositoryRelease{
			Name: "v1.0",
		}, nil)
		client.On("DownloadAsset", tMock.Anything).Return([]byte{}, nil)
		newClientFactory := &exd.NewClientFactory{}
		newClientFactory.Add(provider, func(ctx context.Context, httpDoer exd.HTTPDoer) (exd.Client, error) {
			return client, nil
		})
		exd.NewClientRegistry = newClientFactory

		manifester := &mock.Manifester{}
		manifester.On("Load", tMock.Anything).Return(&exd.Manifest{}, nil)
		manifester.On("Flush", tMock.Anything, tMock.Anything).Return(nil)
		defer manifester.AssertCalled(m.T(), "Flush", tMock.Anything, tMock.Anything)

		installer := &mock.Installer{}
		installer.On("Prepare", tMock.Anything).Return(nil)
		installer.On("Install", tMock.Anything, tMock.Anything).Return(nil)

		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		manager, err := exd.NewManager(ctx, httpDoer, manifester, installer)
		if err != nil {
			panic(err)
		}

		remotePath := "gojek/optimus-extension-valor"
		commandName := "valor"

		actualErr := manager.Install(remotePath, commandName)

		m.NoError(actualErr)
	})
}

func TestManager(t *testing.T) {
	suite.Run(t, &ManagerTestSuite{})
}

func TestNewManager(t *testing.T) {
	t.Run("should return nil and error if context is nil", func(t *testing.T) {
		var ctx context.Context
		httpDoer := &mock.HTTPDoer{}
		manifester := &mock.Manifester{}
		installer := &mock.Installer{}

		actualManager, actualErr := exd.NewManager(ctx, httpDoer, manifester, installer)

		assert.Nil(t, actualManager)
		assert.Error(t, actualErr)
	})

	t.Run("should return nil and error if http doer is nil", func(t *testing.T) {
		ctx := context.Background()
		var httpDoer exd.HTTPDoer
		manifester := &mock.Manifester{}
		installer := &mock.Installer{}

		actualManager, actualErr := exd.NewManager(ctx, httpDoer, manifester, installer)

		assert.Nil(t, actualManager)
		assert.Error(t, actualErr)
	})

	t.Run("should return nil and error if manifester is nil", func(t *testing.T) {
		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		var manifester exd.Manifester
		installer := &mock.Installer{}

		actualManager, actualErr := exd.NewManager(ctx, httpDoer, manifester, installer)

		assert.Nil(t, actualManager)
		assert.Error(t, actualErr)
	})

	t.Run("should return nil and error if installer is nil", func(t *testing.T) {
		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		manifester := &mock.Manifester{}
		var installer exd.Installer

		actualManager, actualErr := exd.NewManager(ctx, httpDoer, manifester, installer)

		assert.Nil(t, actualManager)
		assert.Error(t, actualErr)
	})

	t.Run("should return manager and nil if no error encountered", func(t *testing.T) {
		ctx := context.Background()
		httpDoer := &mock.HTTPDoer{}
		manifester := &mock.Manifester{}
		installer := &mock.Installer{}

		actualManager, actualErr := exd.NewManager(ctx, httpDoer, manifester, installer)

		assert.NotNil(t, actualManager)
		assert.NoError(t, actualErr)
	})
}
