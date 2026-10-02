package messaging

import (
	"os"

	"github.com/stevanusy21/golang_sandbox/pkg/config"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

const LogLocation = "Order RabbitMQ Connection"

func SetupRabbitMqConnection() *config.RabbitMQ {
	rabbit, err := config.ConnectRabbitMQ(LogLocation, os.Getenv("RABBITMQ_URL"))
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal terhubung ke RabbitMQ", err)
	}

	if err = rabbit.Channel.ExchangeDeclare(
		"payment.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		utils.LogFatal(LogLocation, "Gagal mendeklarasikan exchange", err)
	}

	queue, err := rabbit.Channel.QueueDeclare(
		"order.payment-status.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal mendeklarasikan queue", err)
	}

	if err = rabbit.Channel.QueueBind(
		queue.Name,
		"payment.updated",
		"payment.events",
		false,
		nil,
	); err != nil {
		utils.LogFatal(LogLocation, "Gagal mengikat queue", err)
	}

	return rabbit
}
