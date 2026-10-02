package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/midtrans/midtrans-go/coreapi"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stevanusy21/golang_sandbox/pkg/config"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/gateway"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/repository"
)

const LogLocation = "Payment Usecase"

type PaymentUsecase struct {
	repo    *repository.PaymentRepository
	gateway gateway.PaymentGateway
	rabbit  *config.RabbitMQ
}

type PaymentUpdatedEvent struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}

func NewPaymentUsecase(repo *repository.PaymentRepository, gw gateway.PaymentGateway, rabbit *config.RabbitMQ) *PaymentUsecase {
	return &PaymentUsecase{repo: repo, gateway: gw, rabbit: rabbit}
}

func (u *PaymentUsecase) ProcessPayment(orderId string, amount float64, paymentMethod string) (*domain.PaymentRecord, error) {
	record := &domain.PaymentRecord{
		OrderId:       orderId,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        "PENDING",
	}

	if err := u.repo.SavePayment(record); err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrPaymentProcessFailed, err)
	}

	result, err := u.gateway.Charge(orderId, amount, coreapi.CoreapiPaymentType(paymentMethod))
	if err != nil {
		record.Status = "FAILED"
		if saveErr := u.repo.UpdatePayment(record); saveErr != nil {
			return nil, fmt.Errorf("%w: %v", utils.ErrPaymentProcessFailed, saveErr)
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrPaymentProcessFailed, err)
	}

	transactionId := result.TransactionId
	record.TransactionId = &transactionId
	record.Status = strings.ToUpper(result.Status)

	if err := u.repo.UpdatePayment(record); err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrPaymentProcessFailed, err)
	}

	return record, nil
}

func (u *PaymentUsecase) HandlePaymentStatusCallback(callback domain.MidtransCallbackDTO) error {
	record, err := u.repo.GetPaymentByOrderId(callback.OrderID)
	if err != nil {
		if errors.Is(err, utils.ErrPaymentNotFound) {
			utils.LogWarn(LogLocation, "Callback diterima untuk order yang tidak ditemukan: "+callback.OrderID)
			return nil
		}
		return fmt.Errorf("%w: %v", utils.ErrPaymentCallbackFailed, err)
	}

	record.Status = strings.ToUpper(callback.TransactionStatus)
	if err := u.repo.UpdatePayment(record); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrPaymentCallbackFailed, err)
	}

	if record.Status != "PENDING" {
		if err := u.PublishPaymentUpdatedEvent(record.OrderId, record.Status); err != nil {
			return err
		}
	}
	return nil
}

func (u *PaymentUsecase) PublishPaymentUpdatedEvent(orderId string, status string) error {
	event := PaymentUpdatedEvent{
		OrderId: orderId,
		Status:  status,
	}

	body, err := json.Marshal(event)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrFailedToMarshal.Error(), err)
		return fmt.Errorf("%w: %v", utils.ErrPaymentPublishEventFailed, err)
	}

	err = u.rabbit.Channel.Publish(
		"payment.events",
		"payment.updated",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		utils.LogErrorNoValue(LogLocation, utils.ErrPaymentPublishEventFailed.Error())
		return fmt.Errorf("%w: %v", utils.ErrPaymentPublishEventFailed, err)
	}

	utils.LogInfo(LogLocation, "Berhasil publish event payment.updated ke RabbitMQ untuk Order ID: "+orderId)
	return nil
}
