package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductById(id int) (domain.Product, error) {
	var p domain.Product
	
	query := "SELECT id, name, price FROM products WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return domain.Product{}, err
	}

	return p, nil
}

func (r *ProductRepository) GetAllProducts(filter domain.ProductFilter) ([]domain.Product, error) {
	query := "SELECT id, name, price FROM products WHERE 1=1"

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

	query += filter.BuildOrderBy()

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			log.Println("Error saat scan data:", err)
			continue
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(product *domain.Product) error {
	query := "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRow(query, product.Name, product.Price).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}