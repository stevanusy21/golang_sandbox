package grpc

import (
	"context"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/proto/payment"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/usecase"
)

type PaymentHandler struct {
	payment.UnimplementedPaymentServiceServer
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(usecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: usecase}
}

func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	utils.LogInfo("Payment Handler", fmt.Sprintf("Menerima request pembayaran untuk Order ID: %s dengan jumlah: %.2f via %s", req.OrderId, req.Amount, req.PaymentMethod))

	record, err := h.paymentUsecase.ProcessPayment(req.GetOrderId(), req.GetAmount(), req.GetPaymentMethod())
	if err != nil {
		return &payment.PaymentResponse{
			TransactionId: "",
			IsSuccess:     false,
			Message:       err.Error(),
		}, nil
	}

	return &payment.PaymentResponse{
		TransactionId: *record.TransactionId,
		IsSuccess:     true,
		Message:       "Pembayaran berhasil diproses via " + req.GetPaymentMethod(),
	}, nil
}
