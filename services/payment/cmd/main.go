package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/pkg/config"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/proto/payment"
	deliveryGrpc "github.com/stevanusy21/golang_sandbox/services/payment/internal/delivery/grpc"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/payment/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/gateway"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/usecase"
	"google.golang.org/grpc"
)

const ServiceName = "Payment Service"

func main() {
	err := godotenv.Load("services/payment/.env")
	if err != nil {
		utils.LogErrorNoValue(ServiceName, "File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := config.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal(ServiceName, "Gagal terhubung ke database", err)
	}
	defer db.Close()

	rabbit := setupRabbitMqConnection()
	defer rabbit.Close()

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isProduction := os.Getenv("MIDTRANS_PRODUCTION") == "true"

	midtransGateway := gateway.NewMidtransGateway(serverKey, isProduction)

	paymentRepo := repository.NewPaymentRepository(db)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, midtransGateway, rabbit)
	paymentHandler := deliveryHttp.NewPaymentHandler(paymentUsecase)

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("POST /callback/midtrans", paymentHandler.HandlePaymentStatusCallback)

		httpPort := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

		utils.LogInfo(ServiceName, ServiceName+" (HTTP) sedang berjalan di port "+httpPort)
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			utils.LogFatal(ServiceName, "Gagal menjalankan server HTTP", err)
		}
	}()

	paymentHandlerGrpc := deliveryGrpc.NewPaymentHandler(paymentUsecase)
	grpcServer := grpc.NewServer()
	payment.RegisterPaymentServiceServer(grpcServer, paymentHandlerGrpc)

	grpcPort := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		utils.LogFatal(ServiceName, "Gagal membuka port gRPC", err)
	}

	utils.LogInfo(ServiceName, ServiceName+" (gRPC) sedang berjalan di port "+grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		utils.LogFatal(ServiceName, "Gagal menjalankan server gRPC", err)
	}
}

func setupRabbitMqConnection() *config.RabbitMQ {
	rabbit, err := config.ConnectRabbitMQ(ServiceName, os.Getenv("RABBITMQ_URL"))
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

	return rabbit
}
