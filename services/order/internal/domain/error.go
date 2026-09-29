package domain

import "errors"

var(
	ErrOrderNotFound       = errors.New("Order tidak ditemukan")
	ErrInvalidRequest        = errors.New("Request tidak valid")
	ErrOrderCreationFailed = errors.New("Gagal membuat order")
	ErrOrderUpdateFailed   = errors.New("Gagal mengupdate order")
	ErrOrderDeleteFailed   = errors.New("Gagal menghapus order")
	ErrOrderQueryFailed    = errors.New("Gagal mengambil data order")
	ErrOrderScanFailed     = errors.New("Gagal memproses data order")
)