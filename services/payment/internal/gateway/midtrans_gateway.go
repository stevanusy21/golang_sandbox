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

func (m *MidtransGateway) Charge(orderId string, amount float64, paymentMethod coreapi.CoreapiPaymentType) (*PaymentResult, error) {
	chargeReq := &coreapi.ChargeReq{}

	switch paymentMethod {
	case coreapi.PaymentTypeGopay:
		chargeReq.Gopay = &coreapi.GopayDetails{
			EnableCallback: true,
		}
	default:
		return nil, fmt.Errorf("unsupported payment method: %v", paymentMethod)
	}

	chargeReq.PaymentType = paymentMethod
	chargeReq.TransactionDetails = midtrans.TransactionDetails{
		OrderID:  orderId,
		GrossAmt: int64(amount),
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