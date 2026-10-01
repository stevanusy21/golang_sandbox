package repository

import (
	"database/sql"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
)

const LogLocation = "Order Repository"

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(o *domain.Order) (domain.Order, error) {
	query := `
		INSERT INTO orders (id, user_id, product_id, quantity, total_amount, status) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, product_id, quantity, total_amount, status, created_at, updated_at, deleted_at
	`

	var order domain.Order
	err := r.db.QueryRow(
		query,
		o.Id,
		o.UserId,
		o.ProductId,
		o.Quantity,
		o.TotalAmount,
		o.Status,
	).Scan(
		&order.Id,
		&order.UserId,
		&order.ProductId,
		&order.Quantity,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
	)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrOrderCreationFailed.Error(), err)
		return order, fmt.Errorf("%w: %v", utils.ErrOrderCreationFailed, err)
	}

	return order, nil
}

func (r *OrderRepository) UpdateStatus(orderId string, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := r.db.Exec(query, status, orderId)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrOrderUpdateFailed.Error(), err)
		return fmt.Errorf("%w: %v", utils.ErrOrderUpdateFailed, err)
	}

	return nil
}

func (r *OrderRepository) GetAllOrders(filter domain.OrderFilter) ([]domain.Order, error) {
	qb := utils.NewPaginationQueryBuilder(`
		SELECT 
			id, 
			user_id, 
			product_id, 
			quantity, 
			total_amount, 
			status, 
			created_at,
			updated_at,
			deleted_at
		FROM orders 
		WHERE 1 = 1
	`)

	if filter.Id != "" {
		qb.Where("AND", "id", "=", filter.Id)
	}

	if filter.UserId != 0 {
		qb.Where("AND", "user_id", "=", filter.UserId)
	}

	if filter.ProductId != 0 {
		qb.Where("AND", "product_id", "=", filter.ProductId)
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
		return nil, fmt.Errorf("%w: %v", utils.ErrQueryFailed, err)
	}

	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order

		if err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.ProductId,
			&order.Quantity,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
		); err != nil {
			utils.LogError(LogLocation, utils.ErrScanFailed.Error(), err)
			return nil, fmt.Errorf("%w: %v", utils.ErrScanFailed, err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, fmt.Errorf("%w: %v", utils.ErrQueryFailed, err)
	}

	return orders, nil
}

func (r *OrderRepository) GetOrderById(id string) (domain.Order, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			product_id, 
			quantity, 
			total_amount, 
			status, 
			created_at,
			updated_at,
			deleted_at
		FROM orders 
		WHERE id = $1
	`

	var order domain.Order

	err := r.db.QueryRow(query, id).Scan(
		&order.Id,
		&order.UserId,
		&order.ProductId,
		&order.Quantity,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogErrorNoValue(LogLocation, utils.ErrOrderNotFound.Error())
			return domain.Order{}, utils.ErrOrderNotFound
		}
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return domain.Order{}, fmt.Errorf("%w: %v", utils.ErrQueryFailed, err)
	}

	return order, nil
}
