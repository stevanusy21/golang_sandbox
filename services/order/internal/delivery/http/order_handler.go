package http

import (
	"encoding/json"
	"net/http"

	"github.com/stevanusy21/golang_sandbox/pkg/response"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewOrderHandler(u *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: u,
	}
}

func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var input domain.CheckoutRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	order := &domain.Order{
		Customer: input.Customer,
		TotalAmount: input.TotalAmount,
	}

	err := h.orderUsecase.CreateOrder(order, input.PaymentMethod)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Order berhasil dibuat dan pembayaran diproses",
		"order": order,
	})
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	allowedColumns := map[string]bool{
		"id": true,
		"name": true,
		"price": true,
	}
	
	filter := domain.OrderFilter{
		Id: r.URL.Query().Get("id"),
		Customer: r.URL.Query().Get("customer"),
		TotalAmount: r.URL.Query().Get("total_amount"),
		Status: r.URL.Query().Get("status"),
		Pagination: utils.GeneratePagination(
			r.URL.Query().Get("page"),
			r.URL.Query().Get("limit"),
			r.URL.Query().Get("sort_by"),
			r.URL.Query().Get("sort_dir"),
			allowedColumns,
		),
	}
	
	orders, err := h.orderUsecase.GetAllOrders(filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, orders)
}

