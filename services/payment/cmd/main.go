package main

import (
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/proto/payment"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	deliveryGrpc "github.com/stevanusy21/golang_sandbox/services/payment/internal/delivery/grpc"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/gateway"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/payment/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("services/payment/.env")
	if err != nil {
		utils.LogErrorNoValue("Payment Service", "File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := utils.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal("Payment Service", "Gagal terhubung ke database", err)
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
		utils.LogFatal("Payment Service", "Gagal membuka port", err)
	}

	utils.LogInfo("Payment Service", "Payment Service (gRPC) sedang berjalan di port "+port)

	if err := grpcServer.Serve(listener); err != nil {
		utils.LogFatal("Payment Service", "Gagal menjalankan server gRPC", err)
	}
}
