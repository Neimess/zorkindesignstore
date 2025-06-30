package attribute

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"testing"
	"time"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
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

type PGAttributeRepositorySuite struct {
	suite.Suite
	tc   *testContainer
	repo *PGAttributeRepository
	ctx  context.Context
}

type testContainer struct {
	container testcontainers.Container
	db        *sqlx.DB
}

func (s *PGAttributeRepositorySuite) SetupSuite() {
	t_log.SetDefault(log.New(io.Discard, "", log.LstdFlags))

	ctx := context.Background()
	s.ctx = ctx

	postgresContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute),
		))
	require.NoError(s.T(), err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(s.T(), err)

	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(s.T(), err)

	require.NoError(s.T(), createSchema(db))

	s.tc = &testContainer{
		container: postgresContainer,
		db:        db,
	}
	dep, err := NewDeps(db, slog.New(slog.DiscardHandler))
	require.NoError(s.T(), err)
	s.repo = NewPGAttributeRepository(dep)
}

func (s *PGAttributeRepositorySuite) TearDownSuite() {
	_ = s.tc.db.Close()
	_ = s.tc.container.Terminate(context.Background())
}

func createSchema(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS categories (
		category_id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS attributes (
		attribute_id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		unit VARCHAR(50),
		category_id BIGINT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
		UNIQUE(name, category_id)
	);
	`
	_, err := db.Exec(schema)
	return err
}

func (s *PGAttributeRepositorySuite) createCategory(name string) int64 {
	var id int64
	err := s.tc.db.QueryRow(`INSERT INTO categories (name) VALUES ($1) RETURNING category_id`, name).Scan(&id)
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

	t.Run("duplicate name", func(t *testing.T) {
		attrs := []attr.Attribute{
			{Name: "X", Unit: ptr("см"), CategoryID: catID},
			{Name: "X", Unit: ptr("см"), CategoryID: catID},
		}
		err := s.repo.SaveBatch(s.ctx, attrs)
		assert.Error(t, err)
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
