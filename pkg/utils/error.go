package utils

import "errors"

var (
	//Umum
	ErrInvalidRequest        = errors.New("Kesalahan pada request")
	ErrInternalServerError   = errors.New("Terjadi kesalahan pada server")
	ErrQueryFailed           = errors.New("Gagal melakukan query")
	ErrScanFailed			 = errors.New("Gagal memproses data dari query database")

	//User
	ErrUserNotFound          = errors.New("User tidak ditemukan")
	ErrUsernameAlreadyExists = errors.New("Username sudah terdaftar")
	ErrEmailAlreadyExists    = errors.New("Email sudah terdaftar")
	ErrHashPassword          = errors.New("Gagal melakukan hash password")
	ErrUserCreationFailed    = errors.New("Gagal menyimpan user")
	ErrUserUpdateFailed      = errors.New("Gagal update user")
	ErrUserDeleteFailed      = errors.New("Gagal menghapus user")
	
	//Order
	ErrOrderNotFound       = errors.New("Order tidak ditemukan")
	ErrOrderCreationFailed = errors.New("Gagal membuat order")
	ErrOrderUpdateFailed   = errors.New("Gagal mengupdate order")
	ErrOrderDeleteFailed   = errors.New("Gagal menghapus order")

	//Payment
	ErrPaymentNotFound       = errors.New("Payment tidak ditemukan")
	ErrPaymentProcessFailed  = errors.New("Gagal memproses payment")
	ErrPaymentSaveFailed     = errors.New("Gagal menyimpan payment")
	ErrPaymentUpdateFailed   = errors.New("Gagal mengupdate payment")
	ErrPaymentPublishEventFailed = errors.New("Gagal publish event")
	ErrPaymentCallbackFailed = errors.New("Gagal callback payment")

	//Product
	ErrProductNotFound       = errors.New("Produk tidak ditemukan")
	ErrProductCreationFailed = errors.New("Gagal membuat produk")
	ErrProductUpdateFailed   = errors.New("Gagal mengupdate produk")
	ErrProductDeleteFailed   = errors.New("Gagal menghapus produk")
	ErrProductStockNotEnough = errors.New("Stok produk tidak mencukupi")
)
