package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/api/proto/payment"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/config"
	deliveryGrpc "github.com/stevanusy21/golang_sandbox/services/payment/internal/delivery/grpc"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/gateway"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("services/payment/.env")
	if err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isProduction := os.Getenv("MIDTRANS_PRODUCTION") == "true"

	midtransGateway := gateway.NewMidtransGateway(serverKey, isProduction)
	
	paymentRepo := repository.NewPaymentRepository(db)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, midtransGateway)
	paymentHandler := deliveryGrpc.NewPaymentHandler(paymentUsecase)

	grpcServer := grpc.NewServer()

	payment.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	log.Printf("Payment Service (gRPC) sedang berjalan di port %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Gagal menjalankan server gRPC: %v", err)
	}
}
