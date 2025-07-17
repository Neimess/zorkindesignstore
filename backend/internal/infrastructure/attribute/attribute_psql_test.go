package attribute_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"testing"
	"time"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	attrRepo "github.com/Neimess/zorkin-store-project/internal/infrastructure/attribute"
	testsuite "github.com/Neimess/zorkin-store-project/pkg/database/test_suite"
	"github.com/Neimess/zorkin-store-project/pkg/migrator"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	t_log "github.com/testcontainers/testcontainers-go/log"
)

type PGAttributeRepositorySuite struct {
	suite.Suite
	repo *attrRepo.PGAttributeRepository
	ctx  context.Context
	srv  *testsuite.TestServer
	db   *sqlx.DB
}

func (s *PGAttributeRepositorySuite) SetupSuite() {
	t_log.SetDefault(log.New(io.Discard, "", log.LstdFlags))

	srv := testsuite.RunTestServer(s.T())
	require.NotNil(s.T(), srv)

	s.srv = srv
	s.ctx = context.Background()

	require.NoError(s.T(), migrator.Run(srv.Cfg.Storage.DSN(), migrator.Options{Mode: migrator.Up}))

	deps, err := attrRepo.NewDeps(srv.App.DB(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(s.T(), err)

	s.repo = attrRepo.NewPGAttributeRepository(deps)

	s.db = srv.App.DB()
}

func (s *PGAttributeRepositorySuite) TearDownSuite() {
	_ = s.srv.App.DB().Close()
}

func (s *PGAttributeRepositorySuite) createCategory(name string) int64 {
	var id int64
	err := s.db.QueryRow(`INSERT INTO categories (name) VALUES ($1) RETURNING category_id`, name).Scan(&id)
	require.NoError(s.T(), err)
	return id
}

// ─── TESTS ─────────────────────────────────────────────────────────────

func (s *PGAttributeRepositorySuite) Test_SaveBatch() {
	t := s.T()
	catID := s.createCategory(fmt.Sprintf("cat_%d", time.Now().UnixNano()))

	t.Run("ok", func(t *testing.T) {
		attrs := []attr.Attribute{
			{Name: "A", Unit: ptr("кг"), CategoryID: catID},
			{Name: "B", Unit: ptr("шт"), CategoryID: catID},
		}

		err := s.repo.SaveBatch(s.ctx, attrs)
		require.NoError(t, err)
		assert.NotZero(t, attrs[0].ID)
		assert.NotZero(t, attrs[1].ID)
		assert.NotEqual(t, attrs[0].ID, attrs[1].ID)
	})

	t.Run("empty batch", func(t *testing.T) {
		err := s.repo.SaveBatch(s.ctx, []attr.Attribute{})
		assert.NoError(t, err)
	})
}

func (s *PGAttributeRepositorySuite) Test_SaveAndGetByID() {
	catID := s.createCategory(fmt.Sprintf("cat_%d", time.Now().UnixNano()))
	a := &attr.Attribute{Name: "Вес", Unit: ptr("кг"), CategoryID: catID}

	require.NoError(s.T(), s.repo.Save(s.ctx, a))
	assert.NotZero(s.T(), a.ID)

	found, err := s.repo.GetByID(s.ctx, a.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), a.Name, found.Name)
	assert.Equal(s.T(), *a.Unit, *found.Unit)
}

func (s *PGAttributeRepositorySuite) Test_Update() {
	catID := s.createCategory(fmt.Sprintf("cat_%d", time.Now().UnixNano()))
	a := &attr.Attribute{Name: "Initial", Unit: ptr("кг"), CategoryID: catID}

	require.NoError(s.T(), s.repo.Save(s.ctx, a))

	a.Name = "Updated"
	a.Unit = ptr("г")

	updated, err := s.repo.Update(s.ctx, a)
	require.NoError(s.T(), err)
	require.Equal(s.T(), a.Name, updated.Name)
	found, err := s.repo.GetByID(s.ctx, a.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated", found.Name)
	assert.Equal(s.T(), "г", *found.Unit)
}

func (s *PGAttributeRepositorySuite) Test_Delete() {
	catID := s.createCategory(fmt.Sprintf("cat_%d", time.Now().UnixNano()))
	a := &attr.Attribute{Name: "DeleteMe", Unit: ptr("шт"), CategoryID: catID}

	require.NoError(s.T(), s.repo.Save(s.ctx, a))

	require.NoError(s.T(), s.repo.Delete(s.ctx, a.ID))

	_, err := s.repo.GetByID(s.ctx, a.ID)
	assert.Error(s.T(), err)
}

func (s *PGAttributeRepositorySuite) Test_FindByCategory() {
	catID := s.createCategory(fmt.Sprintf("cat_%d", time.Now().UnixNano()))

	attrs := []attr.Attribute{
		{Name: "A", Unit: ptr("кг"), CategoryID: catID},
		{Name: "B", Unit: ptr("шт"), CategoryID: catID},
		{Name: "C", Unit: nil, CategoryID: catID},
	}
	require.NoError(s.T(), s.repo.SaveBatch(s.ctx, attrs))

	found, err := s.repo.FindByCategory(s.ctx, catID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), found, 3)
}

func ptr(s string) *string {
	return &s
}

func TestPGAttributeRepositorySuite(t *testing.T) {
	suite.Run(t, new(PGAttributeRepositorySuite))
}
