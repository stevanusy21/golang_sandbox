package http

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/stevanusy21/golang_sandbox/pkg/request"
	"github.com/stevanusy21/golang_sandbox/pkg/response"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/usecase"
)

type PaymentHandler struct {
	usecase *usecase.PaymentUsecase
}

func NewPaymentHandler(usecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{usecase: usecase}
}

func (h *PaymentHandler) HandlePaymentStatusCallback(w http.ResponseWriter, r *http.Request) {
	callback, err := request.DecodeJSON[domain.MidtransCallbackDTO](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if !verifySignature(callback, serverKey) {
		response.Error(w, http.StatusBadRequest, "Invalid signature")
		return
	}

	if err = h.usecase.HandlePaymentStatusCallback(callback); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "OK",
	})
}

func verifySignature(callback domain.MidtransCallbackDTO, serverKey string) bool {
	raw := callback.OrderID + callback.StatusCode + callback.GrossAmount + serverKey
	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])
	return expected == callback.SignatureKey
}
