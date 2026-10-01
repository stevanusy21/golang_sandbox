package config

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func ConnectRabbitMQ(serviceName, urlStr string) (*RabbitMQ, error) {
	utils.LogInfo(serviceName, "Mencoba konek ke RabbitMQ di "+urlStr+"...")

	conn, err := amqp.Dial(urlStr)
	if err != nil {
		utils.LogFatal(serviceName, "Gagal konek ke RabbitMQ", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		utils.LogFatal(serviceName, "Gagal membuka channel RabbitMQ", err)
		return nil, err
	}

	utils.LogInfo(serviceName, "Berhasil konek ke RabbitMQ")
	return &RabbitMQ{Conn: conn, Channel: ch}, nil
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}
