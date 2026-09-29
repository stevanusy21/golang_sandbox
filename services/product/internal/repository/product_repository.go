package repository

import (
	"database/sql"
	"fmt"
	"log"

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
		WHERE 1=1`

	args := []any{}
	paramIndex := 1

	if filter.ID != "" {
		query += fmt.Sprintf(" AND id = $%d", paramIndex)
		args = append(args, filter.ID)
		paramIndex++
	}

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", paramIndex)
		args = append(args, "%"+filter.Name+"%")
		paramIndex++
	}

	if filter.Price != "" {
		query += fmt.Sprintf(" AND price = $%d", paramIndex)
		args = append(args, filter.Price)
		paramIndex++
	}

	if len(filter.Status) > 0 {
		query += fmt.Sprintf(" AND status = ANY($%d)", paramIndex)
		args = append(args, filter.Status)
		paramIndex++
	}

	query += filter.BuildOrderBy()

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var scanError error = nil

	var products []domain.Product
	for rows.Next() {
		var p domain.Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)

		if err != nil {
			scanError = domain.ErrProductScanFailed
			log.Println("Error saat scan data:", err)
			continue
		}

		products = append(products, p)
	}

	if scanError != nil {
		return nil, scanError
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
