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
	"github.com/stevanusy21/golang_sandbox/proto/product"
	deliveryGrpc "github.com/stevanusy21/golang_sandbox/services/product/internal/delivery/grpc"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/product/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/usecase"
	"google.golang.org/grpc"
)

const LogLocation = "Product Service Main"

func main() {
	//Setup Env
	err := godotenv.Load("services/product/.env")
	if err != nil {
		utils.LogErrorNoValue(LogLocation, "File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	//Setup DB
	db, err := config.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal konek ke Database", err)
	}
	defer db.Close()

	//Setup Repo, Usecase, Handler
	productRepo := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUsecase(productRepo)
	productHandler := deliveryHttp.NewProductHandler(productUseCase)

	//Setup Router
	mux := http.NewServeMux()
	productHandler.RegisterRoutes(mux)

	//Run Server HTTP
	go func() {
		port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
		utils.LogInfo(LogLocation, "Product Service (HTTP) berjalan di port "+port)

		if err := http.ListenAndServe(port, mux); err != nil {
			utils.LogFatal(LogLocation, "Server gagal berjalan", err)
		}
	}()

	//Run Server GRPC
	productHandlerGrpc := deliveryGrpc.NewProductHandlerGrpc(productUseCase)
	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(grpcServer, productHandlerGrpc)

	grpcPort := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal membuka port", err)
	}

	utils.LogInfo(LogLocation, "Product Service (gRPC) berjalan di port "+grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		utils.LogFatal(LogLocation, "Server GRPC gagal berjalan", err)
	}

}
