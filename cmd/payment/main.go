package main

import (
	"log"
	"net"

	"github.com/stevanusy21/golang_sandbox/api/proto/payment"
	deliveryGrpc "github.com/stevanusy21/golang_sandbox/internal/delivery/grpc"
	"google.golang.org/grpc"
)

func main() {
	paymentHandler := deliveryGrpc.NewPaymentHandler()

	grpcServer := grpc.NewServer()

	payment.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	log.Printf("Payment Service (gRPC) sedang berjalan di port %s", port)

	if err := grpcServer.Serve(listener);
	err != nil {
		log.Fatalf("Gagal menjalankan server gRPC: %v", err)
	}
}