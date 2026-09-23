package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *domain.Order) error {
	query := `INSERT INTO orders (id, customer, total_amount, status, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, order.Id, order.Customer, order.TotalAmount, order.Status, order.CreatedAt)

	return err;
}

func (r *OrderRepository) UpdateStatus(orderId string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := r.db.Exec(query, status, orderId)

	return err;
}

func (r *OrderRepository) GetAllOrders(filter domain.OrderFilter) ([]domain.Order, error) {
	query := "SELECT id, customer, total_amount, status, created_at FROM orders WHERE 1 = 1"

	args := []any{}
	paramIndex := 1

	if filter.Id != "" {
		query += fmt.Sprintf(" AND id = $%d", paramIndex)
		args = append(args, filter.Id)
		paramIndex++
	}

	if filter.Customer != "" {
		query += fmt.Sprintf(" AND customer ILIKE $%d", paramIndex)
		args = append(args, filter.Customer)
		paramIndex++
	}

	if filter.TotalAmount != "" {
		query += fmt.Sprintf(" AND total_amount = $%d", paramIndex)
		args = append(args, filter.TotalAmount)
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

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		
		err := rows.Scan(&order.Id, &order.Customer, &order.TotalAmount, &order.Status, &order.CreatedAt)
		if err != nil {
			log.Println("Error saat scan data:", err)
			continue
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}