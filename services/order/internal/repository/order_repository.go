package repository

import (
	"database/sql"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *domain.Order) error {
	query := `
		INSERT INTO orders (id, customer, total_amount, status) 
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, order.Id, order.Customer, order.TotalAmount, order.Status)
	if err != nil {
		utils.LogError("Order Repository", "Error create order", err)
		return domain.ErrOrderCreationFailed
	}

	return nil
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
			customer, 
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

	if filter.Customer != "" {
		qb.Where("AND", "customer", "ILIKE", "%"+filter.Customer+"%")
	}

	if filter.TotalAmount != "" {
		qb.Where("AND", "total_amount", "=", filter.TotalAmount)
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
			&order.Customer,
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
			customer, 
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
		&order.Customer,
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
