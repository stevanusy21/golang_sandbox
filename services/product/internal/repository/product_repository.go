package repository

import (
	"database/sql"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductById(id int) (domain.Product, error) {
	var p domain.Product

	query := `
		SELECT 
			id, 
			name, 
			price, 
			stock, 
			status,
			created_at,
			updated_at,
			deleted_at 
		FROM products 
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)

	return p, err
}

func (r *ProductRepository) GetAllProducts(filter domain.ProductFilter) ([]domain.Product, error) {
	qb := utils.NewPaginationQueryBuilder(`
		SELECT 
			id, 
			name, 
			price, 
			stock, 
			status,
			created_at,
			updated_at,
			deleted_at 
		FROM products 
		WHERE 1=1
	`)

	if filter.Id != "" {
		qb.Where("AND", "id", "=", filter.Id)
	}

	if filter.Name != "" {
		qb.Where("AND", "name", "ILIKE", "%"+filter.Name+"%")
	}

	if filter.Price != "" {
		qb.Where("AND", "price", "=", filter.Price)
	}

	if len(filter.Status) > 0 {
		qb.WhereAny("AND", "status", filter.Status)
	}

	qb.OrderBy(filter.SortBy, filter.SortDir)
	qb.LimitOffset(filter.Limit, filter.Offset)

	query, args := qb.Build()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrProductScanFailed, err)
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, err
}

func (r *ProductRepository) CreateProduct(p *domain.Product) error {
	query := `
		INSERT INTO products (name, price, stock, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, p.Name, p.Price, p.Stock, p.Status)

	return err
}

func (r *ProductRepository) UpdateProduct(id int, p *domain.ProductUpdateRequest) error {
	query := `
		UPDATE products 
		SET 
			name = $1, 
			price = $2, 
			stock = $3, 
			status = $4, 
			updated_at = CURRENT_TIMESTAMP 
		WHERE id = $5 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, p.Name, p.Price, p.Stock, p.Status, id)

	return err
}

func (r *ProductRepository) DeleteProduct(id int) error {
	query := `
		UPDATE products 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(query, id)

	return err
}
