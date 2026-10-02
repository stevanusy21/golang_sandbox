package utils

import "errors"

var (
	//Umum
	ErrInvalidRequest        = errors.New("Kesalahan pada request")
	ErrInternalServerError   = errors.New("Terjadi kesalahan pada server")
	ErrQueryFailed           = errors.New("Gagal melakukan query")
	ErrScanFailed			 = errors.New("Gagal memproses data dari query database")

	//Token
	ErrTokenMissing          = errors.New("Token tidak ditemukan")
	ErrTokenInvalid          = errors.New("Token tidak valid")
	ErrTokenExpired          = errors.New("Token sudah expired")
	ErrTokenCreatedFailed    = errors.New("Gagal membuat token")
	ErrTokenParseFailed    = errors.New("Gagal memparsing token")
	ErrJwtSecretMissing 	= errors.New("JWT Secret tidak ditemukan")

	//User
	ErrUserNotFound          = errors.New("User tidak ditemukan")
	ErrUsernameAlreadyExists = errors.New("Username sudah terdaftar")
	ErrEmailAlreadyExists    = errors.New("Email sudah terdaftar")
	ErrHashPassword          = errors.New("Gagal melakukan hash password")
	ErrUserGetFailed		 = errors.New("Gagal mendapatkan user")
	ErrUserCreationFailed    = errors.New("Gagal menyimpan user")
	ErrUserUpdateFailed      = errors.New("Gagal update user")
	ErrUserDeleteFailed      = errors.New("Gagal menghapus user")
	ErrUserChangePasswordFailed = errors.New("Gagal mengubah password user")
	ErrUserChangeStatusFailed = errors.New("Gagal mengubah status user")
	ErrLoginFailed			 = errors.New("Gagal login")
	
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
