package domain

import (
	"errors"
	"time"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type OrderCreateRequest struct {
	Customer      string  `json:"customer"`
	TotalAmount   float64 `json:"total_amount"`
	PaymentMethod string  `json:"payment_method"`
}

type OrderFilter struct {
	Id          string        `form:"id"`
	Customer    string        `form:"customer"`
	TotalAmount string        `form:"total_amount"`
	Status      []OrderStatus `form:"status"`
	utils.Pagination
}

type OrderDetailResponse struct {
	Id          string      `json:"id"`
	Customer    string      `json:"customer"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

func (c OrderCreateRequest) Validate() error {
	if utils.IsEmpty(c.Customer) {
		return errors.New("Customer tidak boleh kosong")
	}

	if utils.IsEmpty(c.TotalAmount) {
		return errors.New("TotalAmount tidak boleh kosong")
	}

	if utils.IsEmpty(c.PaymentMethod) {
		return errors.New("PaymentMethod tidak boleh kosong")
	}

	if c.TotalAmount < 0 {
		return errors.New("TotalAmount tidak boleh negatif")
	}

	return nil
}
