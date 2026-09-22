package domain

import (
	"time"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type Order struct {
	Id    string     `json:"id"`
	Customer  string  `json:"customer"`
	TotalAmount float64 `json:"total_amount"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CheckoutRequest struct {
	Customer string `json:"customer"`
	TotalAmount float64 `json:"total_amount"`
	PaymentMethod string `json:"payment_method"`
}

type OrderFilter struct {
	Id string
	Customer string
	TotalAmount string
	Status string
	utils.Pagination
}