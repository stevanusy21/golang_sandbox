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
	orderRepo     *repository.OrderRepository
	paymentClient payment.PaymentServiceClient
}

func NewOrderUsecase(repo *repository.OrderRepository, payClient payment.PaymentServiceClient) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:     repo,
		paymentClient: payClient,
	}
}

func (u *OrderUsecase) CreateOrder(payload *domain.OrderCreateRequest) error {
	order := domain.Order{
		Id:          fmt.Sprintf("ORD-%d", time.Now().UnixMilli()),
		Customer:    payload.Customer,
		TotalAmount: payload.TotalAmount,
		Status:      domain.OrderPending,
	}

	err := u.orderRepo.CreateOrder(&order)
	if err != nil {
		return err
	}

	req := &payment.PaymentRequest{
		OrderId:       order.Id,
		Amount:        order.TotalAmount,
		PaymentMethod: payload.PaymentMethod,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	res, err := u.paymentClient.ProcessPayment(ctx, req)
	if err != nil {
		u.orderRepo.UpdateStatus(order.Id, domain.OrderFailed)
		return err
	}

	if !res.IsSuccess {
		u.orderRepo.UpdateStatus(order.Id, domain.OrderFailed)
		return err
	}

	return nil
}

func (u *OrderUsecase) GetAllOrders(filter domain.OrderFilter) ([]domain.OrderDetailResponse, error) {
	orders, err := u.orderRepo.GetAllOrders(filter)
	if err != nil {
		return nil, err
	}

	response := make([]domain.OrderDetailResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, order.ToOrderDetailResponse())
	}

	return response, nil
}

func (u *OrderUsecase) GetOrderById(id string) (domain.OrderDetailResponse, error) {
	order, err := u.orderRepo.GetOrderById(id)
	if err != nil {
		return domain.OrderDetailResponse{}, err
	}

	return order.ToOrderDetailResponse(), nil
}
