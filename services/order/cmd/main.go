package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/proto/payment"
	"github.com/stevanusy21/golang_sandbox/proto/product"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
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

	// Koneksi ke Payment Service
	grpcConnPayment, err := utils.ConnectGrpc("Order Service", os.Getenv("GRPC_PAYMENT_ADDRESS"))
	if err != nil {
		utils.LogFatal("Order Service", "Gagal konek ke Payment Service", err)
	}
	defer grpcConnPayment.Close()

	paymentClient := payment.NewPaymentServiceClient(grpcConnPayment)

	// Koneksi ke Product Service
	grpcConnProduct, err := utils.ConnectGrpc("Order Service", os.Getenv("GRPC_PRODUCT_ADDRESS"))
	if err != nil {
		utils.LogFatal("Order Service", "Gagal konek ke Product Service", err)
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

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	utils.LogInfo("Order Service", "Order Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal("Order Service", "Server gagal berjalan", err)
	}
}
