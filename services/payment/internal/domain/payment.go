package domain

import "time"

type PaymentRecord struct {
	Id              int        `json:"id"`
	OrderId         string     `json:"order_id"`
	Amount          float64    `json:"amount"`
	PaymentMethod   string     `json:"payment_method"`
	Status          string     `json:"status"`
	TransactionId   *string    `json:"transaction_id"`
	CreatedAt		time.Time  `json:"created_at"`
	UpdatedAt		time.Time  `json:"updated_at"`
}
