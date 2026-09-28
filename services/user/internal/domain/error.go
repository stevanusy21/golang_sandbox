package domain

import "errors"

var (
	ErrUserNotFound          = errors.New("User tidak ditemukan")
	ErrUsernameAlreadyExists = errors.New("Username sudah terdaftar")
	ErrEmailAlreadyExists    = errors.New("Email sudah terdaftar")
	ErrHashPassword          = errors.New("Gagal melakukan hash password")
	ErrInvalidRequest        = errors.New("Kesalahan pada request")
	ErrInternalServerError   = errors.New("Terjadi kesalahan pada server")
	ErrUserCreationFailed    = errors.New("Gagal menyimpan user")
	ErrUserUpdateFailed      = errors.New("Gagal update user")
	ErrUserDeleteFailed      = errors.New("Gagal menghapus user")
)