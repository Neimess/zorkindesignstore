package category

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"sort"
	"testing"
	"time"

	cat "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/pkg/app_error"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	t_log "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName = "testdb"
	dbUser = "testuser"
	dbPass = "testpass"
)

type CategoryRepositorySuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sqlx.DB
	repo      *PGCategoryRepository
	ctx       context.Context
}

func (s *CategoryRepositorySuite) SetupSuite() {
	t_log.SetDefault(log.New(io.Discard, "", log.LstdFlags))

	s.ctx = context.Background()

	postgresContainer, err := postgres.Run(s.ctx, "postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
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

	err = s.createSchema()
	require.NoError(s.T(), err)

	s.repo = NewPGCategoryRepository(Deps{db: s.db, log: slog.New(slog.DiscardHandler)})
}

func (s *CategoryRepositorySuite) SetupTest() {
	_, _ = s.db.Exec(`DELETE FROM categories`)
}

func (s *CategoryRepositorySuite) TearDownSuite() {
	_ = s.db.Close()
	_ = s.container.Terminate(s.ctx)
}

func (s *CategoryRepositorySuite) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS categories (
		category_id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

func (s *CategoryRepositorySuite) createCategory(name string) int64 {
	catObj := &cat.Category{Name: name}
	created, err := s.repo.Create(s.ctx, catObj)
	require.NoError(s.T(), err)
	return created.ID
}

// ─── Tests ─────────────────────────────────────────────────────────────

func (s *CategoryRepositorySuite) Test_CreateAndGetByID() {
	id := s.createCategory(fmt.Sprintf("cat_%d", time.Now().UnixNano()))

	fetched, err := s.repo.GetByID(s.ctx, id)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), id, fetched.ID)
}

func (s *CategoryRepositorySuite) Test_GetByID_NotFound() {
	_, err := s.repo.GetByID(s.ctx, 99999)
	assert.ErrorIs(s.T(), err, app_error.ErrNotFound)
}

func (s *CategoryRepositorySuite) Test_Update() {
	id := s.createCategory("oldname")

	err := s.repo.Update(s.ctx, id, "newname")
	require.NoError(s.T(), err)

	fetched, err := s.repo.GetByID(s.ctx, id)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "newname", fetched.Name)
}

func (s *CategoryRepositorySuite) Test_Update_NotFound() {
	err := s.repo.Update(s.ctx, 99999, "name")
	assert.ErrorIs(s.T(), err, app_error.ErrNotFound)
}

func (s *CategoryRepositorySuite) Test_Delete() {
	id := s.createCategory("todelete")

	err := s.repo.Delete(s.ctx, id)
	require.NoError(s.T(), err)

	_, err = s.repo.GetByID(s.ctx, id)
	assert.ErrorIs(s.T(), err, app_error.ErrNotFound)
}

func (s *CategoryRepositorySuite) Test_Delete_NotFound() {
	err := s.repo.Delete(s.ctx, 99999)
	assert.ErrorIs(s.T(), err, app_error.ErrNotFound)
}

func (s *CategoryRepositorySuite) Test_List() {
	names := []string{"bbb", "aaa", "ccc"}
	for _, n := range names {
		s.createCategory(n)
	}

	listed, err := s.repo.List(s.ctx)
	require.NoError(s.T(), err)
	sort.Slice(listed, func(i, j int) bool {
		return listed[i].Name < listed[j].Name
	})
	assert.Equal(s.T(), []string{"aaa", "bbb", "ccc"}, []string{listed[0].Name, listed[1].Name, listed[2].Name})
}

func TestCategoryRepositorySuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositorySuite))
}
