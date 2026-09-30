package domain

import (
	"errors"
	"time"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type OrderCreateRequest struct {
	UserId        int        `json:"user_id"`
	ProductId     int        `json:"product_id"`
	Quantity      int        `json:"quantity"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type OrderFilter struct {
	Id        string        `form:"id"`
	UserId    int           `form:"user_id"`
	ProductId int           `form:"product_id"`
	Status    []OrderStatus `form:"status"`
	utils.Pagination
}

type OrderDetailResponse struct {
	Id          string      `json:"id"`
	UserId      int         `json:"user_id"`
	ProductId   int         `json:"product_id"`
	Quantity    int         `json:"quantity"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

func (c OrderCreateRequest) Validate() error {
	if c.UserId <= 0 {
		return errors.New("UserId yang valid diperlukan")
	}

	if c.ProductId <= 0 {
		return errors.New("ProductId yang valid diperlukan")
	}

	if !c.PaymentMethod.IsValid() {
		return errors.New("PaymentMethod yang valid diperlukan")
	}

	if c.Quantity <= 0 {
		return errors.New("Quantity minimal 1")
	}

	return nil
}
