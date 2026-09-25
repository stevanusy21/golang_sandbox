package repository

import (
	"database/sql"

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