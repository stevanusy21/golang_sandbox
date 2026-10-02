package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stevanusy21/golang_sandbox/pkg/config"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/proto/payment"
	"github.com/stevanusy21/golang_sandbox/proto/product"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
)

const LogLocation = "Order Service Main"

func main() {
	//Setup Env
	err := godotenv.Load("services/order/.env")
	if err != nil {
		utils.LogErrorNoValue(LogLocation, "File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	//Setup DB
	db, err := config.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal konek ke Database", err)
	}
	defer db.Close()

	// Setup koneksi Payment Service
	grpcConnPayment, err := config.ConnectGrpc(LogLocation, os.Getenv("GRPC_PAYMENT_ADDRESS"))
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal konek ke Payment Service", err)
	}
	defer grpcConnPayment.Close()

	paymentClient := payment.NewPaymentServiceClient(grpcConnPayment)

	//Setup koneksi Product Service
	grpcConnProduct, err := config.ConnectGrpc(LogLocation, os.Getenv("GRPC_PRODUCT_ADDRESS"))
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal konek ke Product Service", err)
	}
	defer grpcConnProduct.Close()

	productClient := product.NewProductServiceClient(grpcConnProduct)

	//Setup Repo, Usecase, Handler
	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, paymentClient, productClient)
	orderHandler := deliveryHttp.NewOrderHandler(orderUsecase)

	//Setup Router
	mux := http.NewServeMux()
	orderHandler.RegisterRoutes(mux)

	//Setup Rabbit MQ
	rabbit := setupRabbitMqConnection()
	defer rabbit.Close()

	//Run Consumer RabbitMq 
	go startPaymentStatusConsumer(rabbit.Channel, orderUsecase)

	//Run Server HTTP
	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	utils.LogInfo(LogLocation, "Order Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal(LogLocation, "Server gagal berjalan", err)
	}
}

func setupRabbitMqConnection() *config.RabbitMQ {
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

type PaymentUpdatedEvent struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}

func startPaymentStatusConsumer(ch *amqp.Channel, orderUsecase *usecase.OrderUsecase) {
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
