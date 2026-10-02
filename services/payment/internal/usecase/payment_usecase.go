package usecase

import (
	"encoding/json"
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
		return nil, err
	}

	result, err := u.gateway.Charge(orderId, amount, coreapi.CoreapiPaymentType(paymentMethod))
	if err != nil {
		record.Status = "FAILED"
		if updateErr := u.repo.UpdatePayment(record); updateErr != nil {
			return nil, err
		}
		return nil, err
	}

	transactionId := result.TransactionId
	record.TransactionId = &transactionId
	record.Status = strings.ToUpper(result.Status)

	if err := u.repo.UpdatePayment(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (u *PaymentUsecase) HandlePaymentStatusCallback(callback domain.MidtransCallbackDTO) error {
	record, err := u.repo.GetPaymentByOrderId(callback.OrderID)
	if err != nil {
		return err
	}

	record.Status = strings.ToUpper(callback.TransactionStatus)
	if err := u.repo.UpdatePayment(record); err != nil {
		return err
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
		return utils.ErrFailedToMarshal
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
		utils.LogError(LogLocation, utils.ErrPaymentPublishEventFailed.Error(), err)
		return utils.ErrPaymentPublishEventFailed
	}

	utils.LogInfo(LogLocation, "Berhasil publish event payment.updated ke RabbitMQ untuk Order ID: "+orderId)
	return nil
}
