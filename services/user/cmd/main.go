package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stevanusy21/golang_sandbox/pkg/config"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	deliveryHttp "github.com/stevanusy21/golang_sandbox/services/user/internal/delivery/http"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/repository"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/usecase"
)

func main() {
	err := godotenv.Load("services/user/.env")
	if err != nil {
		utils.LogErrorNoValue("User Service", "File Env tidak ditemukan, menggunakan environment variables sistem")
	}

	db, err := config.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal("User Service", "Gagal terhubung ke database", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := deliveryHttp.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	utils.LogInfo("User Service", "User Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal("User Service", "Server gagal berjalan", err)
	}
}
