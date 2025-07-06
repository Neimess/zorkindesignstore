package product

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	database "github.com/Neimess/zorkin-store-project/pkg/database"
	tx "github.com/Neimess/zorkin-store-project/pkg/database/tx"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	serviceDom "github.com/Neimess/zorkin-store-project/internal/domain/service"
)

// Deps holds dependencies for PGProductRepository
type Deps struct {
	db  *sqlx.DB
	log *slog.Logger
}

// NewDeps constructs Deps or returns an error if missing params
func NewDeps(db *sqlx.DB, log *slog.Logger) (Deps, error) {
	if db == nil {
		return Deps{}, fmt.Errorf("product repository: missing database connection")
	}
	if log == nil {
		return Deps{}, fmt.Errorf("product repository: missing logger")
	}
	return Deps{db: db, log: log.With("component", "PGProductRepository")}, nil
}

// PGProductRepository implements CRUD for products
type PGProductRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

// NewPGProductRepository creates a new repository instance
func NewPGProductRepository(d Deps) *PGProductRepository {
	return &PGProductRepository{db: d.db, log: d.log}
}

// Create inserts a product and sets its ID and CreatedAt
func (r *PGProductRepository) Create(ctx context.Context, p *prodDom.Product) (*prodDom.Product, error) {
	const query = `INSERT INTO products(name, price, description, category_id, image_url)
		VALUES ($1,$2,$3,$4,$5) RETURNING product_id, created_at`

	var id int64
	var created time.Time
	err := database.WithQuery(ctx, r.log, query, func() error {
		return r.db.QueryRowContext(ctx, query,
			p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL,
		).Scan(&id, &created)
	})
	if err := r.mapPostgreSQLError(err); err != nil {
		return nil, err
	}
	p.ID, p.CreatedAt = id, created
	return p, nil
}

// CreateWithAttrs runs product and attribute insert in one transaction
func (r *PGProductRepository) CreateWithAttrs(ctx context.Context, p *prodDom.Product) (*prodDom.Product, error) {
	return tx.RunInTx(ctx, r.db, func(tx *sqlx.Tx) (*prodDom.Product, error) {
		id, err := r.insertProductTx(ctx, tx, p)
		if err != nil {
			return nil, err
		}
		if len(p.Attributes) > 0 {
			if err := r.insertAttributesTx(ctx, tx, id, p.Attributes); err != nil {
				return nil, err
			}
		}
		if len(p.Services) > 0 {
			if err := r.insertServicesTx(ctx, tx, id, p.Services); err != nil {
				return nil, err
			}
		}
		p.ID = id
		return p, nil
	})
}

// GetWithAttrs retrieves a product and its attributes
func (r *PGProductRepository) GetWithAttrs(ctx context.Context, id int64) (*prodDom.Product, error) {
	prod, err := r.fetchProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	attrs, err := r.fetchAttributes(ctx, id)
	if err != nil {
		return nil, err
	}
	prod.Attributes = attrs
	return prod, nil
}

// ListByCategory lists all products in a category with attributes
func (r *PGProductRepository) ListByCategory(ctx context.Context, catID int64) ([]prodDom.Product, error) {
	const query = `
		SELECT product_id, name, price, description, category_id, image_url, created_at
		FROM products WHERE category_id = $1
	`

	var raws []productRow
	err := r.withQuery(ctx, query, func() error {
		return r.db.SelectContext(ctx, &raws, query, catID)
	})
	if err := r.mapPostgreSQLError(err); err != nil {
		return nil, err
	}

	// сначала конвертируем raws -> prods
	prods := make([]prodDom.Product, 0, len(raws))
	for _, raw := range raws {
		prods = append(prods, *raw.toDomain(nil))
	}

	// если пусто — возвращаем [] и nil
	if len(prods) == 0 {
		return prods, nil
	}

	// собираем атрибуты пачкой
	ids := make([]int64, len(prods))
	for i := range prods {
		ids[i] = prods[i].ID
	}

	rows, err := r.fetchAttributesBatch(ctx, ids)
	if err != nil {
		return nil, err
	}

	// мапим атрибуты по ID
	byID := make(map[int64]*prodDom.Product, len(prods))
	for i := range prods {
		byID[prods[i].ID] = &prods[i]
	}

	for _, ar := range rows {
		pa := ar.toDomainAttr()
		byID[ar.ProductID].Attributes = append(byID[ar.ProductID].Attributes, pa)
	}

	return prods, nil
}

// UpdateWithAttrs updates product fields and replaces attributes
func (r *PGProductRepository) UpdateWithAttrs(ctx context.Context, p *prodDom.Product) (*prodDom.Product, error) {
	return tx.RunInTx(ctx, r.db, func(tx *sqlx.Tx) (*prodDom.Product, error) {
		const upd = `
			UPDATE products
			SET name = $1,
				price = $2,
				description = $3,
				category_id = $4,
				image_url = $5
			WHERE product_id = $6
		`
		res, err := tx.ExecContext(ctx, upd, p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL, p.ID)
		if err != nil {
			return nil, r.mapPostgreSQLError(err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rows == 0 {
			return nil, prodDom.ErrProductNotFound
		}

		// Чистим старые атрибуты
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM product_attributes WHERE product_id = $1`, p.ID); err != nil {
			return nil, r.mapPostgreSQLError(err)
		}
		// Чистим старые услуги
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM product_services WHERE product_id = $1`, p.ID); err != nil {
			return nil, r.mapPostgreSQLError(err)
		}
		// Если есть новые атрибуты — вставляем
		if len(p.Attributes) > 0 {
			ids := make([]int64, len(p.Attributes))
			aids := make([]int64, len(p.Attributes))
			vals := make([]string, len(p.Attributes))
			for i, a := range p.Attributes {
				ids[i], aids[i], vals[i] = p.ID, a.AttributeID, a.Value
			}

			const ins = `
				INSERT INTO product_attributes (product_id, attribute_id, value)
				SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::text[])
			`
			if _, err := tx.ExecContext(ctx, ins,
				pq.Array(ids), pq.Array(aids), pq.Array(vals)); err != nil {
				return nil, r.mapPostgreSQLError(err)
			}
		}
		// Если есть новые услуги — вставляем
		if len(p.Services) > 0 {
			if err := r.insertServicesTx(ctx, tx, p.ID, p.Services); err != nil {
				return nil, err
			}
		}
		return p, nil
	})
}

// Delete removes product; cascades attributes via FK
func (r *PGProductRepository) Delete(ctx context.Context, id int64) error {
	const del = `DELETE FROM products WHERE product_id = $1`
	err := r.withQuery(ctx, del, func() error {
		res, err := r.db.ExecContext(ctx, del, id)
		if err != nil {
			return err
		}
		cnt, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if cnt == 0 {
			return prodDom.ErrProductNotFound
		}
		return nil
	})
	return r.mapPostgreSQLError(err)
}

// --- Internal helpers below ---

func (r *PGProductRepository) insertProductTx(ctx context.Context, tx *sqlx.Tx, p *prodDom.Product) (int64, error) {
	const q = `INSERT INTO products(name, price, description, category_id, image_url) VALUES($1,$2,$3,$4,$5) RETURNING product_id`
	var id int64
	if err := r.withQuery(ctx, q, func() error {
		return tx.QueryRowContext(ctx, q, p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL).Scan(&id)
	}); err != nil {
		return 0, r.mapPostgreSQLError(err)
	}
	return id, nil
}

func (r *PGProductRepository) insertAttributesTx(
	ctx context.Context,
	tx *sqlx.Tx,
	prodID int64,
	attrs []prodDom.ProductAttribute,
) error {
	if len(attrs) == 0 {
		return nil
	}

	prodIDs := make([]int64, len(attrs))
	attrIDs := make([]int64, len(attrs))
	values := make([]string, len(attrs))

	for i, a := range attrs {
		prodIDs[i] = prodID
		attrIDs[i] = a.AttributeID
		values[i] = a.Value
	}

	const query = `
		INSERT INTO product_attributes (product_id, attribute_id, value)
		SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::text[])
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		pq.Array(prodIDs),
		pq.Array(attrIDs),
		pq.Array(values),
	)
	if err != nil {
		return r.mapPostgreSQLError(err)
	}

	return nil
}

func (r *PGProductRepository) insertServicesTx(ctx context.Context, tx *sqlx.Tx, prodID int64, services []serviceDom.Service) error {
	if len(services) == 0 {
		return nil
	}
	prodIDs := make([]int64, len(services))
	svcIDs := make([]int64, len(services))
	for i, s := range services {
		prodIDs[i] = prodID
		svcIDs[i] = s.ID
	}
	const query = `INSERT INTO product_services (product_id, service_id) SELECT * FROM UNNEST($1::bigint[], $2::bigint[]) ON CONFLICT DO NOTHING`
	_, err := tx.ExecContext(ctx, query, pq.Array(prodIDs), pq.Array(svcIDs))
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	return nil
}

func (r *PGProductRepository) fetchProduct(ctx context.Context, id int64) (*prodDom.Product, error) {
	const q = `SELECT product_id, name, price, description, category_id, image_url, created_at FROM products WHERE product_id=$1`
	var raw productRow
	err := r.withQuery(ctx, q, func() error {
		return r.db.GetContext(ctx, &raw, q, id)
	})
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	prod := raw.toDomain(nil)
	// Получаем услуги
	services, err := r.fetchServices(ctx, id)
	if err != nil {
		return nil, err
	}
	prod.Services = services
	return prod, nil
}

func (r *PGProductRepository) fetchAttributes(ctx context.Context, pid int64) ([]prodDom.ProductAttribute, error) {
	const q = `SELECT pa.attribute_id, pa.value, a.name, a.unit FROM product_attributes pa JOIN attributes a ON a.attribute_id=pa.attribute_id WHERE pa.product_id=$1`
	var rows []productAttributeRow
	err := r.db.SelectContext(ctx, &rows, q, pid)
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	var res []prodDom.ProductAttribute
	for _, row := range rows {
		res = append(res, row.toDomainAttr())
	}
	return res, nil
}

func (r *PGProductRepository) fetchAttributesBatch(ctx context.Context, ids []int64) ([]productAttributeRow, error) {
	const q = `SELECT pa.product_id, pa.attribute_id, pa.value, a.name, a.unit FROM product_attributes pa JOIN attributes a USING(attribute_id) WHERE pa.product_id=ANY($1)`
	var rows []productAttributeRow
	err := r.withQuery(ctx, q, func() error {
		return r.db.SelectContext(ctx, &rows, q, pq.Array(ids))
	})
	if err := r.mapPostgreSQLError(err); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *PGProductRepository) fetchServices(ctx context.Context, prodID int64) ([]serviceDom.Service, error) {
	const q = `SELECT s.service_id, s.name, s.description, s.price FROM services s JOIN product_services ps ON s.service_id = ps.service_id WHERE ps.product_id = $1`
	var rows []struct {
		ID          int64   `db:"service_id"`
		Name        string  `db:"name"`
		Description *string `db:"description"`
		Price       float64 `db:"price"`
	}
	err := r.db.SelectContext(ctx, &rows, q, prodID)
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	services := make([]serviceDom.Service, 0, len(rows))
	for _, row := range rows {
		services = append(services, serviceDom.Service{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Price:       row.Price,
		})
	}
	return services, nil
}

func (r *PGProductRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGProductRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}
