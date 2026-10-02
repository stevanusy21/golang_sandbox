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

const LogLocation = "User Service Main"

func main() {
	//Setup Env
	err := godotenv.Load("services/user/.env")
	if err != nil {
		utils.LogErrorNoValue(LogLocation, "File Env tidak ditemukan, menggunakan environment variables sistem")
	}

	//Setup DB
	db, err := config.ConnectDB("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		utils.LogFatal(LogLocation, "Gagal terhubung ke database", err)
	}
	defer db.Close()

	//Setup Repo, Usecase, Handler
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := deliveryHttp.NewUserHandler(userUsecase)

	//Setup Router
	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	//Run Server HTTP
	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	utils.LogInfo(LogLocation, "User Service (HTTP) berjalan di port "+port)

	if err := http.ListenAndServe(port, mux); err != nil {
		utils.LogFatal(LogLocation, "Server gagal berjalan", err)
	}
}
