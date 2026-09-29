package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/user/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/usecase"
)

func main() {
	err := godotenv.Load("services/user/.env")
	if err != nil {
		log.Println("File Env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := utils.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := deliveryHttp.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("PATCH /users/{id}/profile", userHandler.UpdateUser)
	mux.HandleFunc("PUT /users/{id}/password", userHandler.ChangePasswordUser)
	mux.HandleFunc("PATCH /users/{id}/status", userHandler.ChangeStatusUser)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserById)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	fmt.Printf("User Service (HTTP) berjalan di port %s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
