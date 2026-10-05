package repository

import (
	"database/sql"

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
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return domain.Order{}, utils.ErrQueryFailed
	}

	return order, nil
}

func (r *OrderRepository) UpdateStatus(orderId string, status domain.OrderStatus) error {
	query := `
		UPDATE orders 
		SET status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, status, orderId)
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
		return utils.ErrOrderNotFound
	}

	return nil
}

func (r *OrderRepository) GetAllOrders(filter domain.OrderFilter) ([]domain.Order, int, error) {
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

	applyOrderFilters(qb, filter)
	qb.OrderBy(filter.SortBy, filter.SortDir)
	qb.LimitOffset(filter.Limit, filter.Offset)
	dataQuery, dataArgs := qb.Build()

	count := utils.NewPaginationQueryBuilder(`
		SELECT COUNT(*)
		FROM orders
		WHERE 1 = 1
	`)
	applyOrderFilters(count, filter)
	countQuery, countArgs := count.Build()

	var total int
	if err := r.db.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, 0, utils.ErrQueryFailed
	}

	rows, err := r.db.Query(dataQuery, dataArgs...)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, 0, utils.ErrQueryFailed
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
			return nil, 0, utils.ErrScanFailed
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, 0, utils.ErrQueryFailed
	}

	return orders, total, nil
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
		WHERE id = $1 AND deleted_at IS NULL
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
			return domain.Order{}, utils.ErrOrderNotFound
		}
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return domain.Order{}, utils.ErrQueryFailed
	}
	return order, nil
}

func applyOrderFilters(qb *utils.PaginationQueryBuilder, filter domain.OrderFilter) {
	if filter.Id != "" {
		qb.Where("AND", "id", "ILIKE", "%"+filter.Id+"%")
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
	qb.WhereNull("AND", "deleted_at")
}