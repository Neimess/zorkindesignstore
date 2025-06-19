package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/jmoiron/sqlx"
)

// //go:embed sql/product/get_product_with_characteristics.sql
// var queryGetProductByIdWithCharacteristics string

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *domain.Product, images []domain.ProductImage) (retErr error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if retErr != nil {
			_ = tx.Rollback()
		} else {
			retErr = tx.Commit()
		}
	}()

	const queryProduct = `INSERT INTO products (name, price, category_id) VALUES (?, ?, ?)`
	res, retErr := tx.ExecContext(ctx, queryProduct, p.Name, p.Price, p.CategoryID)
	if retErr != nil {
		return retErr
	}
	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	// Если нет изображений, вставляем пустую запись
	if len(images) > 0 {
		const queryImage = `INSERT INTO product_images (product_id, url, alt_text) VALUES (?, ?, ?)`
		for _, img := range images {
			if retErr != nil {
				break
			}
			_, retErr = tx.ExecContext(ctx, queryImage, p.ID, img.URL, img.AltText)
		}
	} else {
		const queryImage = `INSERT INTO product_images (product_id, url, alt_text) VALUES (?, '', '')`
		_, retErr = tx.ExecContext(ctx, queryImage, p.ID)
	}
	return retErr
}

func (r *ProductRepository) GetProduct(ctx context.Context, id int64) (*domain.ProductWithImages, error) {
	var p struct {
		domain.Product
		CategoryName string `db:"category_name"`
	}
	const query = `
		SELECT product_id, name, price, category_id, category_name
		FROM products 
		JOIN USING (category_id)
		WHERE product_id = ?`
	if err := r.db.GetContext(ctx, &p, query, id); err != nil {
		return nil, err
	}

	var imgs []domain.ProductImage
	if err := r.db.SelectContext(ctx, &imgs, `SELECT image_id, product_id, url, alt_text FROM product_images WHERE product_id = ? ORDER BY image_id`, p.ID); err != nil {
		return nil, err
	}

	return &domain.ProductWithImages{
		Product: p.Product,
		Images:  imgs,
	}, nil
}

// GetProductByID возвращает продукт по ID с его изображениями.
func (r *ProductRepository) GetByID(ctx context.Context, productID int64) (*domain.ProductWithCharacteristics, error) {
	var p domain.ProductWithCharacteristics

	const qProduct = `
		SELECT product_id, name, price, category_id
		FROM products
		WHERE product_id = ?
	`
	if err := r.db.GetContext(ctx, &p, qProduct, productID); err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	const qImages = `
		SELECT image_id, url, alt_text
		FROM product_images
		WHERE product_id = ?
	`
	if err := r.db.SelectContext(ctx, &p.Images, qImages, productID); err != nil {
		return nil, fmt.Errorf("get images: %w", err)
	}

	const qChars = `
		SELECT characteristic_id, characteristic_name
		FROM category_characteristics
		WHERE category_id = ?
	`
	if err := r.db.SelectContext(ctx, &p.Characteristics, qChars, p.CategoryID); err != nil {
		return nil, fmt.Errorf("get characteristics: %w", err)
	}

	return &p, nil
}

// GetAllProducts возвращает все продукты с их изображениями.
func (r *ProductRepository) GetAllProductsInCategory(ctx context.Context, category_id int64) ([]domain.ProductWithImages, error) {
	var products []domain.Product
	if err := r.db.SelectContext(ctx, &products, `SELECT product_id, name, price, category_id FROM products WHERE category_id = ?`); err != nil {
		return nil, err
	}
	result := make([]domain.ProductWithImages, 0, len(products))
	for _, p := range products {
		var imgs []domain.ProductImage
		if err := r.db.SelectContext(ctx, &imgs, `SELECT image_id, product_id, url, alt_text FROM product_images WHERE product_id = ? ORDER BY image_id`, p.ID); err != nil {
			return nil, err
		}
		result = append(result, domain.ProductWithImages{Product: p, Images: imgs})
	}
	return result, nil
}
