package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/proto/payment"
	"github.com/stevanusy21/golang_sandbox/proto/product"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
)

const LogLocation = "Order Usecase"

type OrderUsecase struct {
	orderRepo     *repository.OrderRepository
	paymentClient payment.PaymentServiceClient
	productClient product.ProductServiceClient
}

func NewOrderUsecase(repo *repository.OrderRepository, paymentClient payment.PaymentServiceClient, productClient product.ProductServiceClient) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:     repo,
		paymentClient: paymentClient,
		productClient: productClient,
	}
}

func (u *OrderUsecase) CreateOrder(payload *domain.OrderCreateRequest) (domain.OrderDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prod, err := u.productClient.GetProductDetail(ctx, &product.ProductDetailRequest{
		ProductId: int32(payload.ProductId),
	})
	if err != nil {
		utils.LogError(LogLocation, "Error get product detail", err)
		return domain.OrderDetailResponse{}, fmt.Errorf("%w: %v", utils.ErrOrderCreationFailed, err)
	}

	deductRes, err := u.productClient.DeductStock(ctx, &product.DeductStockRequest{
		ProductId: int32(payload.ProductId),
		Quantity:  int32(payload.Quantity),
	})

	if err != nil || !deductRes.IsSuccess {
		return domain.OrderDetailResponse{}, fmt.Errorf("%w: %v", utils.ErrOrderCreationFailed, err)
	}

	order := domain.Order{
		Id:          fmt.Sprintf("ORD-%d", time.Now().UnixMilli()),
		UserId:      payload.UserId,
		ProductId:   payload.ProductId,
		Quantity:    payload.Quantity,
		TotalAmount: float64(prod.Price) * float64(payload.Quantity),
		Status:      domain.OrderPending,
	}

	createdOrder, err := u.orderRepo.CreateOrder(&order)
	if err != nil {
		fmt.Println("Error di create order")
		return domain.OrderDetailResponse{}, err
	}

	go u.processPaymentBackground(
		createdOrder.Id,
		createdOrder.TotalAmount,
		payload.PaymentMethod,
	)

	return createdOrder.ToOrderDetailResponse(), nil
}

func (u *OrderUsecase) processPaymentBackground(orderId string, amount float64, paymentMethod domain.PaymentMethod) {
	defer func() {
		if err := recover(); err != nil {
			utils.LogErrorNoValue("Order Background", fmt.Sprintf("Panic recovered in background: %v", err))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &payment.PaymentRequest{
		OrderId:       orderId,
		Amount:        amount,
		PaymentMethod: string(paymentMethod),
	}

	res, err := u.paymentClient.ProcessPayment(ctx, req)
	if err != nil || !res.IsSuccess {
		utils.LogError("Order Background", fmt.Sprintf("Payment FAILED for Order %s", orderId), err)
		u.orderRepo.UpdateStatus(orderId, domain.OrderFailed)
		return
	}

	utils.LogInfo("Order Background", fmt.Sprintf("Payment CREATED for Order %s", orderId))
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

func (u *OrderUsecase) UpdateOrderStatus(orderId string, status string) error {
	parsedStatus, err := domain.ParseOrderStatus(status)
	if err != nil {
		return err
	}
	return u.orderRepo.UpdateStatus(orderId, parsedStatus)
}
