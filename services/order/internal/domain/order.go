package domain

import (
	"time"
)

type Order struct {
	Id          string
	Customer    string
	TotalAmount float64
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (o Order) ToOrderDetailResponse() OrderDetailResponse {
	return OrderDetailResponse{
		Id:          o.Id,
		Customer:    o.Customer,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
	}
}
