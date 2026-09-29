package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/api/proto/payment"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load("services/order/.env")
	if err != nil {
		utils.LogErrorNoValue("Order Service", "File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := utils.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal("Order Service", "Gagal konek ke Database", err)
	}
	defer db.Close()

	utils.LogInfo("Order Service", "Mencoba konek ke Payment Service (gRPC) di localhost:50051...")
	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.LogFatal("Order Service", "Gagal konek ke Payment Service", err)
	}
	defer grpcConn.Close()

	paymentClient := payment.NewPaymentServiceClient(grpcConn)

	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, paymentClient)
	orderHandler := deliveryHttp.NewOrderHandler(orderUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /checkout", orderHandler.Checkout)
	mux.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrderById)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	utils.LogInfo("Order Service", "Order Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal("Order Service", "Server gagal berjalan", err)
	}
}
