package product_test

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

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	serviceDom "github.com/Neimess/zorkin-store-project/internal/domain/service"
	prodRepo "github.com/Neimess/zorkin-store-project/internal/infrastructure/product"
)

type PGProductRepositorySuite struct {
	suite.Suite
	repo *prodRepo.PGProductRepository
	ctx  context.Context
	srv  *testsuite.TestServer
	db   *sqlx.DB
}

func (s *PGProductRepositorySuite) SetupSuite() {
	log.SetOutput(io.Discard)

	srv := testsuite.RunTestServer(s.T())
	require.NotNil(s.T(), srv)

	s.srv = srv
	s.ctx = context.Background()

	require.NoError(s.T(), migrator.Run(srv.Cfg.Storage.DSN(), migrator.Options{Mode: migrator.Up}))

	deps, err := prodRepo.NewDeps(srv.App.DB(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(s.T(), err)

	s.repo = prodRepo.NewPGProductRepository(deps)

	s.db = srv.App.DB()
}

func (s *PGProductRepositorySuite) TearDownSuite() {
	_ = s.srv.App.DB().Close()
}

func (s *PGProductRepositorySuite) createCategory(name string) int64 {
	uniqueName := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO categories(name) VALUES ($1) RETURNING category_id`,
		uniqueName,
	).Scan(&id)
	require.NoError(s.T(), err)
	return id
}

func (s *PGProductRepositorySuite) createAttribute(name, unit string, categoryID int64) int64 {
	uniqueName := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO attributes(name, unit, category_id) VALUES ($1, $2, $3) RETURNING attribute_id`,
		uniqueName, unit, categoryID,
	).Scan(&id)
	require.NoError(s.T(), err)
	return id
}

func ptr(s string) *string { return &s }

func (s *PGProductRepositorySuite) Test_ProductWithServicesAndAttributes() {
	catID := s.createCategory("cat")
	serviceID := s.createService("delivery", 100)
	attrID := s.createAttribute("weight", "kg", catID)

	p := &prodDom.Product{
		Name:       "TestProduct",
		Price:      123.45,
		CategoryID: catID,
		Attributes: []prodDom.ProductAttribute{{
			AttributeID: attrID,
			Value:       "10",
			Attribute: attrDom.Attribute{
				ID:         attrID,
				Name:       "weight",
				Unit:       ptr("kg"),
				CategoryID: catID,
			},
		}},
		Services: []serviceDom.Service{{
			ID:    serviceID,
			Name:  "delivery",
			Price: 100,
		}},
	}
	created, err := s.repo.CreateWithAttrs(s.ctx, p)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), created.ID)
}

func (s *PGProductRepositorySuite) Test_ProductWithNoServicesOrAttributes() {
	catID := s.createCategory("cat")
	p := &prodDom.Product{
		Name:       "NoLinks",
		Price:      50,
		CategoryID: catID,
	}
	created, err := s.repo.Create(s.ctx, p)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), created.ID)
}

func (s *PGProductRepositorySuite) Test_UpdateProductWithAttrsAndServices() {
	catID := s.createCategory("cat")
	serviceID := s.createService("install", 200)
	attrID := s.createAttribute("height", "cm", catID)
	p := &prodDom.Product{
		Name:       "ToUpdate",
		Price:      77.7,
		CategoryID: catID,
	}
	created, err := s.repo.Create(s.ctx, p)
	require.NoError(s.T(), err)
	created.Name = "UpdatedName"
	created.Attributes = []prodDom.ProductAttribute{{
		AttributeID: attrID,
		Value:       "180",
		Attribute: attrDom.Attribute{
			ID:         attrID,
			Name:       "height",
			Unit:       ptr("cm"),
			CategoryID: catID,
		},
	}}
	created.Services = []serviceDom.Service{{
		ID:    serviceID,
		Name:  "install",
		Price: 200,
	}}
	_, err = s.repo.UpdateWithAttrs(s.ctx, created)
	require.NoError(s.T(), err)
	_, err = s.repo.Get(s.ctx, created.ID)
	require.NoError(s.T(), err)
}

func TestPGProductRepositorySuite(t *testing.T) {
	suite.Run(t, new(PGProductRepositorySuite))
}

func (s *PGProductRepositorySuite) createService(name string, price float64) int64 {
	uniqueName := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO services(name, price) VALUES ($1, $2) RETURNING service_id`,
		uniqueName, price,
	).Scan(&id)
	require.NoError(s.T(), err)
	return id
}
