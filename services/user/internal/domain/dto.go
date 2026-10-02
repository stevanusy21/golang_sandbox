package domain

import (
	"errors"
	"time"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type UserFilter struct {
	ID       int       `form:"id"`
	Username string       `form:"username"`
	Email    string       `form:"email"`
	Status   []UserStatus `form:"status"`
	utils.Pagination
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type UserChangePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserChangeStatusRequest struct {
	Status UserStatus `json:"status"`
}

type UserDetailResponse struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (r UserCreateRequest) Validate() error {
	if utils.IsEmpty(r.Username) {
		return errors.New("Username tidak boleh kosong")
	}

	if utils.IsEmpty(r.Email) {
		return errors.New("Email tidak boleh kosong")
	}

	if utils.IsEmpty(r.Password) {
		return errors.New("Password tidak boleh kosong")
	}

	if utils.MinLength(r.Password, 8) {
		return errors.New("Password harus minimal 8 karakter")
	}

	if !utils.IsEmail(r.Email) {
		return errors.New("Format email tidak valid")
	}

	return nil
}

func (r UserChangePasswordRequest) Validate() error {
	if utils.IsEmpty(r.OldPassword) {
		return errors.New("Password lama tidak boleh kosong")
	}

	if utils.IsEmpty(r.NewPassword) {
		return errors.New("Password baru tidak boleh kosong")
	}

	if utils.IsEmpty(r.ConfirmPassword) {
		return errors.New("Confirmasi password tidak boleh kosong")
	}

	if utils.MinLength(r.NewPassword, 8) {
		return errors.New("Password baru harus minimal 8 karakter")
	}

	if r.NewPassword != r.ConfirmPassword {
		return errors.New("Password baru dan konfirmasi password tidak sama")
	}

	return nil
}

func (r UserChangeStatusRequest) Validate() error {
	if utils.IsEmpty(string(r.Status)) {
		return errors.New("Status tidak boleh kosong")
	}

	if !r.Status.IsValid() {
		return errors.New("Status tidak valid")
	}

	return nil
}

func (r UserUpdateRequest) Validate() error {
	if r.Username == nil && r.Email == nil {
		return errors.New("Tidak ada data profil yang diubah")
	}

	if r.Email != nil && !utils.IsEmail(*r.Email) {
		return errors.New("Format email tidak valid")
	}

	return nil
}

func (r LoginRequest) Validate() error {
	if utils.IsEmpty(r.Email) {
		return errors.New("Email tidak boleh kosong")
	}

	if utils.IsEmpty(r.Password) {
		return errors.New("Password tidak boleh kosong")
	}

	if !utils.IsEmail(r.Email) {
		return errors.New("Format email tidak valid")
	}

	return nil
}
