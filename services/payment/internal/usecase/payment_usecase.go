package usecase

import (
	"fmt"
	"strings"

	"github.com/stevanusy21/golang_sandbox/services/payment/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/gateway"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/repository"
)

type PaymentUsecase struct {
	repo    *repository.PaymentRepository
	gateway gateway.PaymentGateway
}

func NewPaymentUsecase(repo *repository.PaymentRepository, gw gateway.PaymentGateway) *PaymentUsecase {
	return &PaymentUsecase{repo: repo, gateway: gw}
}

func (u *PaymentUsecase) ProcessPayment(orderId string, amount float64, paymentMethod string) (*domain.PaymentRecord, error) {
	record := &domain.PaymentRecord{
		OrderId:       orderId,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        "PENDING",
	}

	result, err := u.gateway.Charge(orderId, amount, paymentMethod)
	if err != nil {
		record.Status = "FAILURE"
		if saveErr := u.repo.SavePayment(record); saveErr != nil {
			fmt.Printf("CRITICAL: Gagal menyimpan record pembayaran yang gagal: %v", saveErr)
		}
		return nil, fmt.Errorf("Gagal memproses pembayaran: %v", err)
	}

	transactionId := result.TransactionId
	record.TransactionId = &transactionId
	record.Status = strings.ToUpper(result.Status)

	if err := u.repo.SavePayment(record); err != nil {
		return nil, fmt.Errorf("Gagal menyimpan data pembayaran: %v", err)
	}

	return record, nil
}
