package repository

import (
	"database/sql"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *domain.Order) (domain.Order, error) {
	query := `
		INSERT INTO orders (id, user_id, product_id, quantity, total_amount, status) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, product_id, quantity, total_amount, status, created_at, updated_at, deleted_at
	`

	var orderResponse domain.Order
	err := r.db.QueryRow(
		query,
		order.Id,
		order.UserId,
		order.ProductId,
		order.Quantity,
		order.TotalAmount,
		order.Status,
	).Scan(
		&orderResponse.Id,
		&orderResponse.UserId,
		&orderResponse.ProductId,
		&orderResponse.Quantity,
		&orderResponse.TotalAmount,
		&orderResponse.Status,
		&orderResponse.CreatedAt,
		&orderResponse.UpdatedAt,
		&orderResponse.DeletedAt,
	)
	if err != nil {
		utils.LogError("Order Repository", "Error create order", err)
		if err == sql.ErrNoRows {
			return domain.Order{}, domain.ErrOrderNotFound
		}
		return orderResponse, fmt.Errorf("%w: %v", domain.ErrOrderCreationFailed, err)
	}

	return orderResponse, nil
}

func (r *OrderRepository) UpdateStatus(orderId string, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := r.db.Exec(query, status, orderId)
	if err != nil {
		utils.LogError("Order Repository", "Error update status", err)
		return domain.ErrOrderUpdateFailed
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
		utils.LogError("Order Repository", "Error query all orders", err)
		return nil, domain.ErrOrderQueryFailed
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
			utils.LogError("Order Repository", "Error scanning order", err)
			return nil, domain.ErrOrderScanFailed
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		utils.LogError("Order Repository", "Error in rows", err)
		return nil, domain.ErrOrderQueryFailed
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

	if err == sql.ErrNoRows {
		return domain.Order{}, domain.ErrOrderNotFound
	} else if err != nil {
		utils.LogError("Order Repository", "Error get order by id", err)
		return domain.Order{}, domain.ErrOrderQueryFailed
	}

	return order, nil
}
