package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stevanusy21/golang_sandbox/services/order/internal/config"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/order/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/order/internal/usecase"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	productRepo := repository.NewProductRepository(db)

	productUseCase := usecase.NewProductUsecase(productRepo)
	
	myHandler := deliveryHttp.NewHandler(productUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /product", myHandler.GetProductByID)
	mux.HandleFunc("GET /products", myHandler.GetAllProducts)
	mux.HandleFunc("POST /product", myHandler.CreateProduct)

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
