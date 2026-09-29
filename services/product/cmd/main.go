package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/product/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/usecase"
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
	utils.LogInfo("Product Service", "Product Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal("Product Service", "Server gagal berjalan", err)
	}
}
