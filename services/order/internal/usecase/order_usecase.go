package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/stevanusy21/golang_sandbox/api/proto/payment"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
)

type OrderUsecase struct {
	orderRepo *repository.OrderRepository
	paymentClient payment.PaymentServiceClient
}

func NewOrderUsecase(repo *repository.OrderRepository, payClient payment.PaymentServiceClient) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: repo,
		paymentClient: payClient,
	}
}

func (u *OrderUsecase) CreateOrder(order *domain.Order, paymentMethod string) error {
	order.Id = fmt.Sprintf("ORD-%d", time.Now().UnixMilli())
	order.Status = "PENDING"
	order.CreatedAt = time.Now()

	err := u.orderRepo.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("Gagal menyimpan order ke database: %v", err)
	}

	req := &payment.PaymentRequest{
		OrderId: order.Id,
		Amount: order.TotalAmount,
		PaymentMethod: paymentMethod,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	
	defer cancel()

	res, err := u.paymentClient.ProcessPayment(ctx, req)
	if err != nil {
		u.orderRepo.UpdateStatus(order.Id, "FAILED")
		return fmt.Errorf("Gagal menghubungi payment service: %v", err)
	}

	if res.IsSuccess {
		u.orderRepo.UpdateStatus(order.Id, "PAID")
	} else {
		u.orderRepo.UpdateStatus(order.Id, "FAILED")
		return fmt.Errorf("Pembayaran ditolak: %s", res.Message)
	}

	return nil
}

func (u *OrderUsecase) GetAllOrders(filter domain.OrderFilter) ([]domain.Order, error) {
	return u.orderRepo.GetAllOrders(filter)
}

func (u *OrderUsecase) GetOrderById(id string) (domain.Order, error) {
	return u.orderRepo.GetOrderById(id)
}