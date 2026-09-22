package grpc

import (
	"context"
	"log"

	"github.com/stevanusy21/golang_sandbox/api/proto/payment"
)

type PaymentHandler struct {
	payment.UnimplementedPaymentServiceServer
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	log.Printf("Menerima request pembayaran untuk Order ID: %s dengan jumlah: %.2f via %s", 
		req.GetOrderId(), req.GetAmount(), req.GetPaymentMethod())

	res := &payment.PaymentResponse{
		TransactionId: "TRX-" + req.GetOrderId(),
		IsSuccess: true,
		Message: "Pembayaran berhasil diproses via " + req.GetPaymentMethod(),
	}

	return res, nil
}