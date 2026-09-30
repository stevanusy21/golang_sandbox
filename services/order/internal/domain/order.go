package domain

import (
	"time"
)

type Order struct {
	Id          string
	UserId      int
	ProductId   int
	Quantity    int
	TotalAmount float64
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (o Order) ToOrderDetailResponse() OrderDetailResponse {
	return OrderDetailResponse{
		Id:          o.Id,
		UserId:      o.UserId,
		ProductId:   o.ProductId,
		Quantity:    o.Quantity,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
	}
}
