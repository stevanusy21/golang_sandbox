package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/config"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/product/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/usecase"
)

func main() {
	err := godotenv.Load("services/product/.env")
	if err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal konek ke Database: %v", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)

	productUseCase := usecase.NewProductUsecase(productRepo)

	productHandler := deliveryHttp.NewProductHandler(productUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /product", productHandler.GetProductByID)
	mux.HandleFunc("GET /products", productHandler.GetAllProducts)
	mux.HandleFunc("POST /product", productHandler.CreateProduct)
	mux.HandleFunc("PUT /product", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /product", productHandler.DeleteProduct)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	fmt.Printf("Product Service (HTTP) berjalan di port %s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
