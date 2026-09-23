package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/api/proto/payment"
	deliveryGrpc "github.com/stevanusy21/golang_sandbox/services/payment/internal/delivery/grpc"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("services/payment/.env")
	if err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	paymentHandler := deliveryGrpc.NewPaymentHandler()

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
