package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

func ConnectDB(driver string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Gagal melakukan koneksi ke database: %v", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("Gagal melakukan ping ke database: %v", err)
	}

	// Jumlah maksimum koneksi yang terbuka di database (menghindari overload)
	db.SetMaxOpenConns(25)

	// Jumlah koneksi idle (menganggur) yang tetap dipertahankan di memori
	db.SetMaxIdleConns(5)

	// Durasi maksimum koneksi boleh digunakan sebelum dihancurkan dan dibuat baru
	db.SetConnMaxLifetime(5 * time.Minute)

	// Berapa lama koneksi idle boleh bertahan sebelum ditutup otomatis
	db.SetConnMaxIdleTime(2 * time.Minute)

	utils.LogInfo("Database Utils", "Berhasil terhubung ke database")
	return db, nil
}
