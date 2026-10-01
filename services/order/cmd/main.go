package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/proto/payment"
	"github.com/stevanusy21/golang_sandbox/proto/product"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
)

const ServiceName = "Order Service"

func main() {
	err := godotenv.Load("services/order/.env")
	if err != nil {
		utils.LogErrorNoValue(ServiceName, "File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := utils.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal(ServiceName, "Gagal konek ke Database", err)
	}
	defer db.Close()

	// Koneksi ke Payment Service
	grpcConnPayment, err := utils.ConnectGrpc(ServiceName, os.Getenv("GRPC_PAYMENT_ADDRESS"))
	if err != nil {
		utils.LogFatal(ServiceName, "Gagal konek ke Payment Service", err)
	}
	defer grpcConnPayment.Close()

	paymentClient := payment.NewPaymentServiceClient(grpcConnPayment)

	// Koneksi ke Product Service
	grpcConnProduct, err := utils.ConnectGrpc(ServiceName, os.Getenv("GRPC_PRODUCT_ADDRESS"))
	if err != nil {
		utils.LogFatal(ServiceName, "Gagal konek ke Product Service", err)
	}
	defer grpcConnProduct.Close()

	productClient := product.NewProductServiceClient(grpcConnProduct)

	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, paymentClient, productClient)
	orderHandler := deliveryHttp.NewOrderHandler(orderUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /checkout", orderHandler.Checkout)
	mux.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrderById)

	// Koneksi ke Rabbit MQ
	rabbit := setupRabbitMqConnection()
	defer rabbit.Close()

	go startPaymentStatusConsumer(rabbit.Channel, orderUsecase)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	utils.LogInfo("Order Service", "Order Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal("Order Service", "Server gagal berjalan", err)
	}
}

func setupRabbitMqConnection() *utils.RabbitMQ {
	rabbit, err := utils.ConnectRabbitMQ(ServiceName, os.Getenv("RABBITMQ_URL"))
	if err != nil {
		utils.LogFatal(ServiceName, "Gagal terhubung ke RabbitMQ", err)
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
		utils.LogFatal(ServiceName, "Gagal mendeklarasikan exchange", err)
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
		utils.LogFatal(ServiceName, "Gagal mendeklarasikan queue", err)
	}

	if err = rabbit.Channel.QueueBind(
		queue.Name,
		"payment.updated",
		"payment.events",
		false,
		nil,
	); err != nil {
		utils.LogFatal(ServiceName, "Gagal mengikat queue", err)
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
		utils.LogError(ServiceName, "Gagal memulai consumer RabbitMQ", err)
		return
	}

	utils.LogInfo(ServiceName, "Mulai mendengarkan pesan dari RabbitMQ (order.payment-status.queue)...")

	for msg := range msgs {
		var event PaymentUpdatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			utils.LogErrorNoValue(ServiceName, "Gagal unmarshal event RabbitMQ")
			continue
		}

		utils.LogInfo(ServiceName, fmt.Sprintf("Menerima Event Payment: Order ID = %s, Status = %s", event.OrderId, event.Status))

		if err := orderUsecase.UpdateOrderStatus(event.OrderId, event.Status); err != nil {
			utils.LogErrorNoValue(ServiceName, fmt.Sprintf("Gagal mengupdate status order ID %s: %v", event.OrderId, err))
		}
	}
}
