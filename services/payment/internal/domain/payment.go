package domain

import "time"

type PaymentRecord struct {
	Id            int
	OrderId       string
	Amount        float64
	PaymentMethod string
	Status        string
	TransactionId *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
