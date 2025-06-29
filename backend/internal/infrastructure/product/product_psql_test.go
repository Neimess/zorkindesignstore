package product

import (
	"context"
	"io"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	t_log "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/pkg/app_error"
)

const (
	dbName = "testdb"
	dbUser = "testuser"
	dbPass = "testpass"
)

type ProductRepositorySuite struct {
	suite.Suite
	ctx       context.Context
	container testcontainers.Container
	db        *sqlx.DB
	repo      *PGProductRepository
}

func (s *ProductRepositorySuite) SetupSuite() {
	t_log.SetDefault(log.New(io.Discard, "", log.LstdFlags))

	s.ctx = context.Background()

	postgresC, err := postgres.Run(s.ctx, "postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(2*time.Minute),
		),
	)
	require.NoError(s.T(), err)
	s.container = postgresC

	connStr, err := postgresC.ConnectionString(s.ctx, "sslmode=disable")
	require.NoError(s.T(), err)

	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(s.T(), err)
	s.db = db

	err = s.createSchema()
	require.NoError(s.T(), err)

	dep, err := NewDeps(s.db, slog.New(slog.DiscardHandler))
	require.NoError(s.T(), err)
	s.repo = NewPGProductRepository(dep)
}

func (s *ProductRepositorySuite) TearDownSuite() {
	_ = s.db.Close()
	_ = s.container.Terminate(s.ctx)
}

func (s *ProductRepositorySuite) createSchema() error {
	schema := `
CREATE TABLE categories (
  category_id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE
);
CREATE TABLE attributes (
  attribute_id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  unit VARCHAR(50),
  category_id BIGINT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
  UNIQUE(name, category_id)
);
CREATE TABLE products (
  product_id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  price NUMERIC(10,2) NOT NULL,
  description TEXT,
  category_id BIGINT NOT NULL REFERENCES categories(category_id) ON DELETE RESTRICT,
  image_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE product_attributes (
  product_id BIGINT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
  attribute_id BIGINT NOT NULL REFERENCES attributes(attribute_id) ON DELETE RESTRICT,
  value TEXT NOT NULL,
  PRIMARY KEY(product_id, attribute_id)
);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *ProductRepositorySuite) createCategory(name string) int64 {
	var id int64
	err := s.db.QueryRow("INSERT INTO categories(name) VALUES($1) RETURNING category_id", name).Scan(&id)
	s.Require().NoError(err)
	return id
}

func (s *ProductRepositorySuite) createAttribute(name, unit string, catID int64) int64 {
	var id int64
	err := s.db.QueryRow("INSERT INTO attributes(name, unit, category_id) VALUES($1,$2,$3) RETURNING attribute_id", name, unit, catID).Scan(&id)
	s.Require().NoError(err)
	return id
}

func (s *ProductRepositorySuite) Test_CreateGetDelete_ProductLifecycle() {
	catID := s.createCategory("c1")
	attr1 := s.createAttribute("A1", "u1", catID)
	attr2 := s.createAttribute("A2", "u2", catID)

	p := &prodDom.Product{
		Name:        "P1",
		Price:       9.99,
		Description: "desc",
		CategoryID:  catID,
		ImageURL:    "u",
		Attributes: []prodDom.ProductAttribute{
			{AttributeID: attr1, Value: "v1"},
			{AttributeID: attr2, Value: "v2"},
		},
	}

	id, err := s.repo.CreateWithAttrs(s.ctx, p)
	s.Require().NoError(err)
	s.Assert().Equal(id, p.ID)

	got, err := s.repo.GetWithAttrs(s.ctx, id)
	s.Require().NoError(err)
	s.Assert().Equal("P1", got.Name)
	s.Assert().Len(got.Attributes, 2)

	err = s.repo.Delete(s.ctx, id)
	s.Require().NoError(err)

	_, err = s.repo.GetWithAttrs(s.ctx, id)
	s.Assert().ErrorIs(err, app_error.ErrNotFound)
}

func (s *ProductRepositorySuite) Test_ListByCategory_UpdateWithAttrs() {
	catID := s.createCategory("c2")
	list, err := s.repo.ListByCategory(s.ctx, catID)
	s.Require().NoError(err)
	s.Assert().Empty(list)

	attr1 := s.createAttribute("B1", "u", catID)
	p := &prodDom.Product{Name: "PX", Price: 1.23, Description: "d", CategoryID: catID, ImageURL: "i",
		Attributes: []prodDom.ProductAttribute{{AttributeID: attr1, Value: "val"}},
	}
	id, err := s.repo.CreateWithAttrs(s.ctx, p)
	s.Require().NoError(err)

	list, err = s.repo.ListByCategory(s.ctx, catID)
	s.Require().NoError(err)
	s.Assert().Len(list, 1)

	attr2 := s.createAttribute("B2", "u2", catID)
	p.ID = id
	p.Attributes = []prodDom.ProductAttribute{{AttributeID: attr2, Value: "v2"}}
	err = s.repo.UpdateWithAttrs(s.ctx, p)
	s.Require().NoError(err)

	got, err := s.repo.GetWithAttrs(s.ctx, id)
	s.Require().NoError(err)
	s.Assert().Len(got.Attributes, 1)
	s.Assert().Equal(attr2, got.Attributes[0].AttributeID)
}

func TestProductRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProductRepositorySuite))
}
