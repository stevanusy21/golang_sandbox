package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentUpdatedEvent struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}

func StartPaymentStatusConsumer(ch *amqp.Channel, orderUsecase *usecase.OrderUsecase) {
	msgs, err := ch.Consume(
		"order.payment-status.queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		utils.LogError(LogLocation, "Gagal memulai consumer RabbitMQ", err)
		return
	}

	utils.LogInfo(LogLocation, "Mulai mendengarkan pesan dari RabbitMQ (order.payment-status.queue)...")

	for msg := range msgs {
		var event PaymentUpdatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			utils.LogErrorNoValue(LogLocation, utils.ErrFailedToUnmarshal.Error())
			continue
		}

		utils.LogInfo(LogLocation, fmt.Sprintf("Menerima Event Payment: Order ID = %s, Status = %s", event.OrderId, event.Status))

		if err := orderUsecase.UpdateOrderStatus(event.OrderId, event.Status); err != nil {
			utils.LogErrorNoValue(LogLocation, fmt.Sprintf("Gagal mengupdate status order ID %s: %v", event.OrderId, err))
		}
	}
}
