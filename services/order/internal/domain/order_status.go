package domain

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
