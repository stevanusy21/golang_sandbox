package utils

import "errors"

var (
	//Umum
	ErrInvalidRequest      = errors.New("Kesalahan pada request")
	ErrInternalServerError = errors.New("Terjadi kesalahan, mohon coba kembali beberapa saat lagi")
	ErrQueryFailed         = errors.New("Gagal melakukan query")
	ErrScanFailed          = errors.New("Gagal memproses data dari query database")
	ErrGetDataFailed       = errors.New("Gagal mendapatkan data")
	ErrFailedToMarshal     = errors.New("Gagal melakukan marshalling")
	ErrFailedToUnmarshal   = errors.New("Gagal melakukan unmarshalling")

	//Security
	ErrTokenMissing        = errors.New("Token tidak ditemukan")
	ErrTokenInvalid        = errors.New("Token tidak valid")
	ErrTokenExpired        = errors.New("Token sudah expired")
	ErrTokenCreationFailed = errors.New("Gagal membuat token")
	ErrTokenParseFailed    = errors.New("Gagal memparsing token")
	ErrJwtSecretMissing    = errors.New("JWT Secret tidak ditemukan")
	ErrWrongEmailPassword  = errors.New("Email atau password salah")
	ErrUserNotActive       = errors.New("User tidak aktif")
	ErrInvalidSignature	   = errors.New("Signature tidak valid")

	//User
	ErrUserNotFound             = errors.New("User tidak ditemukan")
	ErrUsernameAlreadyExists    = errors.New("Username sudah terdaftar")
	ErrEmailAlreadyExists       = errors.New("Email sudah terdaftar")
	ErrHashPassword             = errors.New("Gagal melakukan hash password")

	//Order
	ErrOrderNotFound       = errors.New("Order tidak ditemukan")

	//Payment
	ErrPaymentNotFound           = errors.New("Payment tidak ditemukan")
	ErrPaymentProcessFailed      = errors.New("Gagal memproses payment")
	ErrPaymentPublishEventFailed = errors.New("Gagal publish event")

	//Midtrans
	ErrMidtrans = errors.New("Terjadi kendala di Midtrans")
	ErrNotSupportedPaymentMethod = errors.New("Metode pembayaran tidak dikenal")

	//Product
	ErrProductNotFound       = errors.New("Produk tidak ditemukan")
	ErrDeductStockFailed     = errors.New("Gagal mengurangi stok")
	ErrAddStockFailed        = errors.New("Gagal menambah stok")
)
