package repository

import (
	"database/sql"
	"errors"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/domain"
)

const LogLocation = "Payment Repository"

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) SavePayment(p *domain.PaymentRecord) error {
	query := `
		INSERT INTO payment_records (order_id, amount, payment_method, status, transaction_id) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(query, p.OrderId, p.Amount, p.PaymentMethod, p.Status, p.TransactionId).Scan(&p.Id)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	return nil
}

func (r *PaymentRepository) UpdatePayment(p *domain.PaymentRecord) error {
	query := `
		UPDATE payment_records 
		SET status = $1, transaction_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE order_id = $3 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, p.Status, p.TransactionId, p.OrderId)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	return nil
}

func (r *PaymentRepository) GetPaymentByOrderId(orderId string) (*domain.PaymentRecord, error) {
	query := `
		SELECT 
			id, 
			order_id, 
			amount, 
			payment_method, 
			status, 
			transaction_id, 
			created_at, 
			updated_at, 
			deleted_at 
		FROM payment_records 
		WHERE order_id = $1 AND deleted_at IS NULL
	`

	var record domain.PaymentRecord
	err := r.db.QueryRow(query, orderId).Scan(
		&record.Id,
		&record.OrderId,
		&record.Amount,
		&record.PaymentMethod,
		&record.Status,
		&record.TransactionId,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrPaymentNotFound
		}
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, utils.ErrQueryFailed
	}

	return &record, nil
}
