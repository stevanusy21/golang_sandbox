package config

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB() *sql.DB {
	dsn := "postgres://postgres:indocyber@localhost:5432/golang"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}

	fmt.Println("Berhasil terhubung ke PostgreSql!")
	return db
}