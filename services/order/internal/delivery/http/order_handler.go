package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/stevanusy21/golang_sandbox/pkg/request"
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
	payload, err := request.DecodeJSON[domain.OrderCreateRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	orderDetailResponse, err := h.orderUsecase.CreateOrder(&payload)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrOrderCreationFailed):
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Order berhasil dibuat dan pembayaran diproses",
		"data":    orderDetailResponse,
	})
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	allowedSortColumns := map[string]bool{
		"id":         true,
		"user_id":    true,
		"product_id": true,
		"quantity":   true,
		"status":     true,
	}

	paramId, err := request.GetStringParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	paramUserId, err := request.GetIntParam(r, "user_id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	paramProductId, err := request.GetIntParam(r, "product_id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	statusQuery := r.URL.Query().Get("status")
	statusSlice := []domain.OrderStatus{}

	if statusQuery != "" {
		for _, s := range strings.Split(statusQuery, ",") {
			if !domain.OrderStatus(s).IsValid() {
				response.Error(w, http.StatusBadRequest, "Status tidak valid")
				return
			}
			statusSlice = append(statusSlice, domain.OrderStatus(s))
		}
	}

	filter := domain.OrderFilter{
		Id:        paramId,
		UserId:    paramUserId,
		ProductId: paramProductId,
		Status:    statusSlice,
		Pagination: utils.GeneratePaginationData(
			r.URL.Query().Get("page"),
			r.URL.Query().Get("limit"),
			r.URL.Query().Get("sort_by"),
			r.URL.Query().Get("sort_dir"),
			allowedSortColumns,
		),
	}

	orders, err := h.orderUsecase.GetAllOrders(filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetStringParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderUsecase.GetOrderById(id)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrOrderNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
			return
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, order)
}
