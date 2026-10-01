package domain

import (
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type OrderStatus string

const (
	OrderPending OrderStatus = "PENDING"
	OrderSuccess OrderStatus = "SUCCESS"
	OrderFailed  OrderStatus = "FAILED"
)

func (o OrderStatus) IsValid() bool {
	switch o {
	case OrderPending, OrderSuccess, OrderFailed:
		return true
	default:
		return false
	}
}

func ParseOrderStatus(status string) (OrderStatus, error) {
	switch status {
	case "PENDING":
		return OrderPending, nil
	case "SUCCESS", "SETTLEMENT":
		return OrderSuccess, nil
	case "FAILED", "CANCEL", "EXPIRE", "FAILURE", "DENY":
		return OrderFailed, nil
	default:
		return "", fmt.Errorf("%w: %v", utils.ErrInvalidRequest, "Status tidak valid")
	}
}
