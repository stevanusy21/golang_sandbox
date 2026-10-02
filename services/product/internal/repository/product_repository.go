package repository

import (
	"database/sql"
	"errors"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/domain"
)

const LogLocation = "Product Repository"

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductById(id int) (domain.Product, error) {
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

	var p domain.Product
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

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, utils.ErrProductNotFound
		}
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return domain.Product{}, utils.ErrQueryFailed
	}

	return p, nil
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
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, utils.ErrQueryFailed
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
			utils.LogError(LogLocation, utils.ErrScanFailed.Error(), err)
			return nil, utils.ErrScanFailed
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, utils.ErrQueryFailed
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(p *domain.Product) error {
	query := `
		INSERT INTO products (name, price, stock, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(query, p.Name, p.Price, p.Stock, p.Status).Scan(&p.ID)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	return nil
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

	result, err := r.db.Exec(query, p.Name, p.Price, p.Stock, p.Status, id)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	if rowsAffected == 0 {
		return utils.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(id int) error {
	query := `
		UPDATE products 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	if rowsAffected == 0 {
		return utils.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) DeductStock(id int, quantity int) error {
	query := `
		UPDATE products 
		SET stock = stock - $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2 AND stock >= $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, quantity, id)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	if rowsAffected == 0 {
		return utils.ErrProductNotFound
	}

	return nil
}

// NotUsedYet
func (r *ProductRepository) AddStock(id int, quantity int) error {
	query := `
		UPDATE products 
		SET stock = stock + $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, quantity, id)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	if rowsAffected == 0 {
		return utils.ErrProductNotFound
	}

	return nil
}
