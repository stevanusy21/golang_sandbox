package domain

import "errors"

var (
	ErrProductNotFound       = errors.New("Produk tidak ditemukan")
	ErrInvalidRequest        = errors.New("Request tidak valid")
	ErrProductCreationFailed = errors.New("Gagal membuat produk")
	ErrProductUpdateFailed   = errors.New("Gagal mengupdate produk")
	ErrProductDeleteFailed   = errors.New("Gagal menghapus produk")
	ErrProductQueryFailed    = errors.New("Gagal mengambil data produk")
	ErrProductScanFailed     = errors.New("Gagal memproses data produk")
)
