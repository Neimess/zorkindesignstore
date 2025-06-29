package preset_test

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

	domPreset "github.com/Neimess/zorkin-store-project/internal/domain/preset"
	presetRepo "github.com/Neimess/zorkin-store-project/internal/infrastructure/preset"
)

const (
	dbName = "testdb"
	dbUser = "testuser"
	dbPass = "testpass"
)

type PresetRepositorySuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sqlx.DB
	repo      *presetRepo.PGPresetRepository
	ctx       context.Context
}

func (s *PresetRepositorySuite) SetupSuite() {
	tlog.SetDefault(log.New(io.Discard, "", log.LstdFlags))

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

	s.repo = presetRepo.NewPGPresetRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))
}

func (s *PresetRepositorySuite) TearDownSuite() {
	_ = s.db.Close()
	_ = s.container.Terminate(s.ctx)
}

func (s *PresetRepositorySuite) createSchema() error {
	schema := `
	CREATE TABLE presets (
		preset_id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		total_price NUMERIC(10,2),
		image_url TEXT,
		created_at TIMESTAMPTZ DEFAULT now()
	);
	CREATE TABLE products (
		product_id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price NUMERIC(10,2) NOT NULL,
		image_url TEXT
	);
	CREATE TABLE preset_items (
		preset_item_id BIGSERIAL PRIMARY KEY,
		preset_id BIGINT NOT NULL REFERENCES presets(preset_id) ON DELETE CASCADE,
		product_id BIGINT NOT NULL REFERENCES products(product_id)
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

func (s *PresetRepositorySuite) createProduct(name string, price float64) int64 {
	var id int64
	err := s.db.QueryRow(`INSERT INTO products (name, price) VALUES ($1, $2) RETURNING product_id`, name, price).Scan(&id)
	require.NoError(s.T(), err)
	return id
}

func (s *PresetRepositorySuite) Test_CreateGetDelete() {
	// Create product first
	prodID := s.createProduct("Sink", 199.99)

	// Create Preset
	in := &domPreset.Preset{
		Name:        "Bathroom Set",
		Description: ptr("Full set for bathroom"),
		TotalPrice:  199.99,
		ImageURL:    ptr("https://example.com/image.jpg"),
		Items: []domPreset.PresetItem{
			{ProductID: prodID},
		},
	}

	id, err := s.repo.Create(s.ctx, in)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), id)

	// Get
	got, err := s.repo.Get(s.ctx, id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Bathroom Set", got.Name)
	require.Len(s.T(), got.Items, 1)
	require.Equal(s.T(), prodID, got.Items[0].ProductID)
	require.Equal(s.T(), "Sink", got.Items[0].Product.Name)

	// Delete
	err = s.repo.Delete(s.ctx, id)
	require.NoError(s.T(), err)

	// Get should fail
	_, err = s.repo.Get(s.ctx, id)
	require.ErrorContains(s.T(), err, "not found")
}

func (s *PresetRepositorySuite) Test_ListDetailedAndShort() {
	prodID := s.createProduct("Toilet", 88.00)

	// Clean up table
	_, _ = s.db.Exec(`DELETE FROM presets`)

	p := &domPreset.Preset{
		Name:        "Toilet Only",
		TotalPrice:  88.0,
		Items:       []domPreset.PresetItem{{ProductID: prodID}},
		Description: ptr("Single toilet"),
	}
	_, err := s.repo.Create(s.ctx, p)
	require.NoError(s.T(), err)

	// ListDetailed
	list, err := s.repo.ListDetailed(s.ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), list)
	require.NotEmpty(s.T(), list[0].Items)

	// ListShort
	short, err := s.repo.ListShort(s.ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), short)
	require.Empty(s.T(), short[0].Items)
}

func ptr(s string) *string { return &s }

func TestPresetRepositorySuite(t *testing.T) {
	suite.Run(t, new(PresetRepositorySuite))
}
