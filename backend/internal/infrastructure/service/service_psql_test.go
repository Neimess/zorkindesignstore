package service_test

import (
	"context"
	"io"
	"log"
	"log/slog"
	"testing"

	testsuite "github.com/Neimess/zorkin-store-project/pkg/database/test_suite"
	"github.com/Neimess/zorkin-store-project/pkg/migrator"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	domService "github.com/Neimess/zorkin-store-project/internal/domain/service"
	serviceRepo "github.com/Neimess/zorkin-store-project/internal/infrastructure/service"
)

func ptr(s string) *string { return &s }

type PGServiceRepositorySuite struct {
	suite.Suite
	repo *serviceRepo.PGServiceRepository
	ctx  context.Context
	srv  *testsuite.TestServer
	db   *sqlx.DB
}

func (s *PGServiceRepositorySuite) SetupSuite() {
	log.SetOutput(io.Discard)

	srv := testsuite.RunTestServer(s.T())
	require.NotNil(s.T(), srv)

	s.srv = srv
	s.ctx = context.Background()

	require.NoError(s.T(), migrator.Run(srv.Cfg.Storage.DSN(), migrator.Options{Mode: migrator.Up}))

	s.repo = serviceRepo.NewPGServiceRepository(srv.App.DB(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	s.db = srv.App.DB()
}

func (s *PGServiceRepositorySuite) TearDownSuite() {
	_ = s.srv.App.DB().Close()
}

func (s *PGServiceRepositorySuite) Test_CRUD() {
	ctx := s.ctx
	svc := &domService.Service{Name: "TestService", Description: ptr("desc"), Price: 10.5}
	created, err := s.repo.Create(ctx, svc)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), created.ID)

	got, err := s.repo.Get(ctx, created.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), svc.Name, got.Name)

	created.Name = "UpdatedService"
	_, err = s.repo.Update(ctx, created)
	require.NoError(s.T(), err)

	got, err = s.repo.Get(ctx, created.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "UpdatedService", got.Name)

	require.NoError(s.T(), s.repo.Delete(ctx, created.ID))
	_, err = s.repo.Get(ctx, created.ID)
	require.Error(s.T(), err)
}

func TestPGServiceRepositorySuite(t *testing.T) {
	suite.Run(t, new(PGServiceRepositorySuite))
}
