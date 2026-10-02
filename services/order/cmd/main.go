package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/pkg/config"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/proto/payment"
	"github.com/stevanusy21/golang_sandbox/proto/product"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/messaging"
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
	rabbit := messaging.SetupRabbitMqConnection()
	defer rabbit.Close()

	//Run Consumer RabbitMq
	go messaging.StartPaymentStatusConsumer(rabbit.Channel, orderUsecase)

	//Run Server HTTP
	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	utils.LogInfo(LogLocation, "Order Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal(LogLocation, "Server gagal berjalan", err)
	}
}
