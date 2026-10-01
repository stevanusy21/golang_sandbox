package utils

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func ConnectRabbitMQ(serviceName, urlStr string) (*RabbitMQ, error) {
	LogInfo(serviceName, "Mencoba konek ke RabbitMQ di " +urlStr+ "...")

	conn, err := amqp.Dial(urlStr)
	if err != nil {
		LogFatal(serviceName, "Gagal konek ke RabbitMQ", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		LogFatal(serviceName, "Gagal membuka channel RabbitMQ", err)
		return nil, err
	}

	LogInfo(serviceName, "Berhasil konek ke RabbitMQ")
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

