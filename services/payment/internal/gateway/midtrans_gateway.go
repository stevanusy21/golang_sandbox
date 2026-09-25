package gateway

import (
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransGateway struct {
	client coreapi.Client
}

func NewMidtransGateway(serverKey string, isProduction bool) *MidtransGateway {
	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}

	client := coreapi.Client{}
	client.New(serverKey, env)

	return &MidtransGateway{
		client: client,
	}
}

func (m *MidtransGateway) Charge(orderId string, amount float64, paymentMethod string) (*PaymentResult, error) {
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeGopay,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: orderId,
			GrossAmt: int64(amount),
		},
		Gopay: &coreapi.GopayDetails{
			EnableCallback: false,
		},
	}

	res, midtransErr := m.client.ChargeTransaction(chargeReq)
	if midtransErr != nil {
		return nil, fmt.Errorf("Midtrans error: %v", midtransErr.GetMessage())
	}

	return &PaymentResult{
		TransactionId: res.TransactionID,
		Status: res.TransactionStatus,
		Message: res.StatusMessage,
	}, nil
}