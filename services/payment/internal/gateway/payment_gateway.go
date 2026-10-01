package gateway

import "github.com/midtrans/midtrans-go/coreapi"

type PaymentResult struct {
	TransactionId string
	Status        string
	Message       string
}

type PaymentGateway interface {
	Charge(orderId string, amount float64, paymentMethod coreapi.CoreapiPaymentType) (*PaymentResult, error)
}
