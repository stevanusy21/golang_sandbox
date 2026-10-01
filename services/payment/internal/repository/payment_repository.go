package repository

import (
	"database/sql"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/services/payment/internal/domain"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) SavePayment(p *domain.PaymentRecord) error {
	query := `INSERT INTO payment_records (order_id, amount, payment_method, status, transaction_id) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, p.OrderId, p.Amount, p.PaymentMethod, p.Status, p.TransactionId)

	return err
}

func (r *PaymentRepository) UpdatePayment(p *domain.PaymentRecord) error {
	query := `
		UPDATE payment_records 
		SET status = $1, transaction_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE order_id = $3
	`

	_, err := r.db.Exec(query, p.Status, p.TransactionId, p.OrderId)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrFailedToUpdatePayment, err)
	}

	return nil
}

func (r *PaymentRepository) GetPaymentByOrderId(orderId string) (*domain.PaymentRecord, error) {
	var record domain.PaymentRecord

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
		WHERE order_id = $1`

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
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", domain.ErrRecordNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrQueryFailed, err)
	}

	return &record, nil
}
