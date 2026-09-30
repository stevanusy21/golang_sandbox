package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/proto/product"
	deliveryGrpc "github.com/stevanusy21/golang_sandbox/services/product/internal/delivery/grpc"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/product/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("services/product/.env")
	if err != nil {
		utils.LogErrorNoValue("Product Service", "File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := utils.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal("Product Service", "Gagal konek ke Database", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUsecase(productRepo)
	productHandler := deliveryHttp.NewProductHandler(productUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", productHandler.GetAllProducts)
	mux.HandleFunc("GET /products/{id}", productHandler.GetProductByID)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	go func() {
		utils.LogInfo("Product Service", "Product Service (HTTP) berjalan di port "+port)

		if err := http.ListenAndServe(port, mux); err != nil {
			utils.LogFatal("Product Service", "Server gagal berjalan", err)
		}
	}()

	productHandlerGrpc := deliveryGrpc.NewProductHandlerGrpc(productUseCase)
	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(grpcServer, productHandlerGrpc)

	grpcPort := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		utils.LogFatal("Product Service", "Gagal membuka port", err)
	}

	utils.LogInfo("Product Service", "Product Service (gRPC) berjalan di port "+grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		utils.LogFatal("Product Service", "Server GRPC gagal berjalan", err)
	}

}
