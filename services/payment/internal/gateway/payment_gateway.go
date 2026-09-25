package gateway

type PaymentResult struct {
	TransactionId string
	Status string
	Message string
}

type PaymentGateway interface {
	Charge(orderId string, amount float64, paymentMethod string) (*PaymentResult, error)
}
