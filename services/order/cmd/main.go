package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stevanusy21/golang_sandbox/api/proto/payment"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/config"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	log.Println("Mencoba konek ke Payment Service (gRPC) di localhost:50051...")
	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Gagal konek ke Payment Service: %v", err)
	}
	defer grpcConn.Close()

	paymentClient := payment.NewPaymentServiceClient(grpcConn)

	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, paymentClient)
	orderHandler := deliveryHttp.NewOrderHandler(orderUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /checkout", orderHandler.Checkout)
	mux.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	mux.HandleFunc("GET /order", orderHandler.GetOrderById)

	port := ":8080"
	fmt.Printf("Order Service (HTTP) berjalan di port %s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
