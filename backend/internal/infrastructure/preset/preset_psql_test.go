package preset_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"testing"
	"time"

	testsuite "github.com/Neimess/zorkin-store-project/pkg/database/test_suite"
	"github.com/Neimess/zorkin-store-project/pkg/migrator"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	domPreset "github.com/Neimess/zorkin-store-project/internal/domain/preset"
	"github.com/Neimess/zorkin-store-project/internal/infrastructure/preset"
)

type PGPresetRepositorySuite struct {
	suite.Suite
	repo *preset.PGPresetRepository
	ctx  context.Context
	srv  *testsuite.TestServer
	db   *sqlx.DB
}

func (s *PGPresetRepositorySuite) SetupSuite() {
	log.SetOutput(io.Discard)

	srv := testsuite.RunTestServer(s.T())
	require.NotNil(s.T(), srv)

	s.srv = srv
	s.ctx = context.Background()

	require.NoError(s.T(), migrator.Run(srv.Cfg.Storage.DSN(), migrator.Options{Mode: migrator.Up}))

	s.repo = preset.NewPGPresetRepository(srv.App.DB(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	s.db = srv.App.DB()
}

func (s *PGPresetRepositorySuite) TearDownSuite() {
	_ = s.srv.App.DB().Close()
}

func (s *PGPresetRepositorySuite) createCategory(name string) int64 {
	uniqueName := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO categories(name) VALUES ($1) RETURNING category_id`,
		uniqueName,
	).Scan(&id)
	require.NoError(s.T(), err)
	return id
}

func (s *PGPresetRepositorySuite) createService(name string, price float64) int64 {
	uniqueName := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO services(name, price) VALUES ($1, $2) RETURNING service_id`,
		uniqueName, price,
	).Scan(&id)
	require.NoError(s.T(), err)
	return id
}

func (s *PGPresetRepositorySuite) createAttribute(name, unit string, categoryID int64) int64 {
	uniqueName := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO attributes(name, unit, category_id) VALUES ($1, $2, $3) RETURNING attribute_id`,
		uniqueName, unit, categoryID,
	).Scan(&id)
	require.NoError(s.T(), err)
	return id
}

func (s *PGPresetRepositorySuite) createProduct(name string, price float64, categoryID int64, attrID int64, svcID int64) int64 {
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO products (name, price, category_id) VALUES ($1, $2, $3) RETURNING product_id`,
		name, price, categoryID,
	).Scan(&id)
	require.NoError(s.T(), err)
	if attrID > 0 {
		_, err = s.db.Exec(`INSERT INTO product_attributes (product_id, attribute_id, value) VALUES ($1, $2, $3)`, id, attrID, "val")
		require.NoError(s.T(), err)
	}
	if svcID > 0 {
		_, err = s.db.Exec(`INSERT INTO product_services (product_id, service_id) VALUES ($1, $2)`, id, svcID)
		require.NoError(s.T(), err)
	}
	return id
}

func ptr(s string) *string { return &s }

func (s *PGPresetRepositorySuite) Test_CreateGetDelete() {
	categoryID := s.createCategory("Bathroom")
	attrID := s.createAttribute("weight", "kg", categoryID)
	svcID := s.createService("delivery", 100)
	prodID := s.createProduct("Sink", 199.99, categoryID, attrID, svcID)
	in := &domPreset.Preset{
		Name:        "Bathroom Set",
		Description: ptr("Full set for bathroom"),
		TotalPrice:  199.99,
		ImageURL:    ptr("https://example.com/image.jpg"),
		Items:       []domPreset.PresetItem{{ProductID: prodID}},
	}
	pRes, err := s.repo.Create(s.ctx, in)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), pRes.ID)
	got, err := s.repo.Get(s.ctx, pRes.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Bathroom Set", got.Name)
	require.Len(s.T(), got.Items, 1)
	require.Equal(s.T(), prodID, got.Items[0].ProductID)
	require.Equal(s.T(), "Sink", got.Items[0].Product.Name)
	require.NoError(s.T(), s.repo.Delete(s.ctx, pRes.ID))
	_, err = s.repo.Get(s.ctx, pRes.ID)
	require.ErrorContains(s.T(), err, "not found")
}

func (s *PGPresetRepositorySuite) Test_ListDetailedAndShort() {
	_, _ = s.db.Exec(`DELETE FROM presets`)
	categoryID := s.createCategory("Bathroom")
	attrID := s.createAttribute("weight", "kg", categoryID)
	svcID := s.createService("delivery", 100)
	prodID := s.createProduct("Sink", 199.99, categoryID, attrID, svcID)
	p := &domPreset.Preset{
		Name:       "Toilet Only",
		TotalPrice: 88.0,
		Items:      []domPreset.PresetItem{{ProductID: prodID}},
	}
	_, err := s.repo.Create(s.ctx, p)
	require.NoError(s.T(), err)
	list, err := s.repo.ListDetailed(s.ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), list)
	require.NotEmpty(s.T(), list[0].Items)
	short, err := s.repo.ListShort(s.ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), short)
	require.Empty(s.T(), short[0].Items)
}

func (s *PGPresetRepositorySuite) Test_UpdatePreset() {
	categoryID := s.createCategory("Bathroom")
	attrA := s.createAttribute("height", "cm", categoryID)
	attrB := s.createAttribute("width", "cm", categoryID)
	svcA := s.createService("install", 200)
	svcB := s.createService("warranty", 50)
	prodA := s.createProduct("Alpha", 50, categoryID, attrA, svcA)
	prodB := s.createProduct("Beta", 75, categoryID, attrB, svcB)
	in := &domPreset.Preset{
		Name:       "Original",
		TotalPrice: 50,
		Items:      []domPreset.PresetItem{{ProductID: prodA}},
	}
	pRes, err := s.repo.Create(s.ctx, in)
	require.NoError(s.T(), err)
	pRes.Name = "Updated"
	pRes.TotalPrice = 75
	pRes.Items = []domPreset.PresetItem{{ProductID: prodB}}
	r, err := s.repo.Update(s.ctx, pRes)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Updated", r.Name)
	require.Equal(s.T(), float64(75), r.TotalPrice)
	got, err := s.repo.Get(s.ctx, pRes.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Updated", got.Name)
	require.Len(s.T(), got.Items, 1)
	require.Equal(s.T(), prodB, got.Items[0].ProductID)
}

func (s *PGPresetRepositorySuite) Test_UpdateIdempotent() {
	categoryID := s.createCategory("Bathroom")
	attrID := s.createAttribute("depth", "cm", categoryID)
	svcID := s.createService("support", 30)
	prod := s.createProduct("Gamma", 20, categoryID, attrID, svcID)
	p := &domPreset.Preset{Name: "Same", TotalPrice: 20, Items: []domPreset.PresetItem{{ProductID: prod}}}
	pRes, err := s.repo.Create(s.ctx, p)
	require.NoError(s.T(), err)
	_, err = s.repo.Update(s.ctx, pRes)
	require.NoError(s.T(), err)
	_, err = s.repo.Update(s.ctx, pRes)
	require.NoError(s.T(), err)
}

func TestPGPresetRepositorySuite(t *testing.T) {
	suite.Run(t, new(PGPresetRepositorySuite))
}
