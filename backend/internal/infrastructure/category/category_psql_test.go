package category_test

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
	categoryRepo "github.com/Neimess/zorkin-store-project/internal/infrastructure/category"
	"github.com/Neimess/zorkin-store-project/pkg/app_error"
	testsuite "github.com/Neimess/zorkin-store-project/pkg/database/test_suite"
	"github.com/Neimess/zorkin-store-project/pkg/migrator"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	t_log "github.com/testcontainers/testcontainers-go/log"
)

type CategoryRepositorySuite struct {
	suite.Suite
	db   *sqlx.DB
	srv  *testsuite.TestServer
	repo *categoryRepo.PGCategoryRepository
	ctx  context.Context
}

func (s *CategoryRepositorySuite) SetupSuite() {
	t_log.SetDefault(log.New(io.Discard, "", log.LstdFlags))

	srv := testsuite.RunTestServer(s.T())
	require.NotNil(s.T(), srv)

	s.srv = srv
	s.ctx = context.Background()

	require.NoError(s.T(), migrator.Run(srv.Cfg.Storage.DSN(), migrator.Options{Mode: migrator.Up}))

	deps, err := categoryRepo.NewDeps(srv.App.DB(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(s.T(), err)

	s.repo = categoryRepo.NewPGCategoryRepository(deps)

	s.db = srv.App.DB()
}

func (s *CategoryRepositorySuite) SetupTest() {
	_, _ = s.db.Exec(`DELETE FROM categories`)
}

func (s *CategoryRepositorySuite) TearDownSuite() {
	_ = s.db.Close()
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
	_, err := s.repo.GetByID(s.ctx, 30000)
	assert.ErrorIs(s.T(), err, app_error.ErrNotFound)
}

func (s *CategoryRepositorySuite) Test_Update() {
	id := s.createCategory("oldname")
	categ := &cat.Category{ID: id, Name: "newname"}
	updated, err := s.repo.Update(s.ctx, categ)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "newname", updated.Name)
	fetched, err := s.repo.GetByID(s.ctx, id)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "newname", fetched.Name)
}

func (s *CategoryRepositorySuite) Test_Update_NotFound() {
	categ := &cat.Category{ID: 30000, Name: "newname"}
	_, err := s.repo.Update(s.ctx, categ)
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
	err := s.repo.Delete(s.ctx, 30000)
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
