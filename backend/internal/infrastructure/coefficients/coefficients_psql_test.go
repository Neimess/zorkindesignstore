package coefficients_test

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tlog "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"
	repo "github.com/Neimess/zorkin-store-project/internal/infrastructure/coefficients"
)

type CoefficientRepositorySuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sqlx.DB
	repo      *repo.PGCoefficientsRepository
	ctx       context.Context
}

func (s *CoefficientRepositorySuite) SetupSuite() {
	tlog.SetDefault(log.New(io.Discard, "", log.LstdFlags))

	s.ctx = context.Background()
	postgresContainer, err := postgres.Run(s.ctx, "postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute),
		),
	)
	require.NoError(s.T(), err)
	s.container = postgresContainer

	connStr, err := postgresContainer.ConnectionString(s.ctx, "sslmode=disable")
	require.NoError(s.T(), err)

	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(s.T(), err)
	s.db = db

	require.NoError(s.T(), s.createSchema())
	s.repo = repo.NewPGCoefficientsRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))
}

func (s *CoefficientRepositorySuite) TearDownSuite() {
	_ = s.db.Close()
	_ = s.container.Terminate(s.ctx)
}

func (s *CoefficientRepositorySuite) createSchema() error {
	schema := `
	CREATE TABLE coefficients (
		coefficient_id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		value NUMERIC(10, 4) NOT NULL
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

func (s *CoefficientRepositorySuite) Test_CreateGetDelete() {
	in := &domCoeff.Coefficient{
		Name:  "TestCoeff",
		Value: 1.2345,
	}
	created, err := s.repo.Create(s.ctx, in)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), created.ID)

	got, err := s.repo.Get(s.ctx, created.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "TestCoeff", got.Name)
	require.Equal(s.T(), 1.2345, got.Value)

	require.NoError(s.T(), s.repo.Delete(s.ctx, created.ID))
	_, err = s.repo.Get(s.ctx, created.ID)
	require.Error(s.T(), err)
}

func (s *CoefficientRepositorySuite) Test_List() {
	_, _ = s.db.Exec(`DELETE FROM coefficients`)
	_, _ = s.repo.Create(s.ctx, &domCoeff.Coefficient{Name: "A", Value: 1.1})
	_, _ = s.repo.Create(s.ctx, &domCoeff.Coefficient{Name: "B", Value: 2.2})
	list, err := s.repo.List(s.ctx)
	require.NoError(s.T(), err)
	require.Len(s.T(), list, 2)
}

func (s *CoefficientRepositorySuite) Test_Update() {
	in := &domCoeff.Coefficient{Name: "ToUpdate", Value: 3.3}
	created, err := s.repo.Create(s.ctx, in)
	require.NoError(s.T(), err)
	created.Name = "Updated"
	created.Value = 4.4
	updated, err := s.repo.Update(s.ctx, created)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Updated", updated.Name)
	require.Equal(s.T(), 4.4, updated.Value)
	got, err := s.repo.Get(s.ctx, created.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Updated", got.Name)
	require.Equal(s.T(), 4.4, got.Value)
}

func TestCoefficientRepositorySuite(t *testing.T) {
	suite.Run(t, new(CoefficientRepositorySuite))
}
