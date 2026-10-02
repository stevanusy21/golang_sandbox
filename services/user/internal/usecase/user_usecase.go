package usecase

import (
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/token"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/repository"
)

type UserUsecase struct {
	repo *repository.UserRepository
}

func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) CreateUser(payload *domain.UserCreateRequest) (int, error) {
	existingUser, err := u.repo.CheckExistsBy("username", payload.Username, nil)
	if err == nil && existingUser {
		return 0, utils.ErrUsernameAlreadyExists
	}

	existingUser, err = u.repo.CheckExistsBy("email", payload.Email, nil)
	if err == nil && existingUser {
		return 0, utils.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return 0, utils.ErrHashPassword
	}

	user := domain.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashedPassword),
		Status:   domain.UserActive,
	}

	if err := u.repo.CreateUser(&user); err != nil {
		return 0, fmt.Errorf("%w: %v", utils.ErrUserCreationFailed, err)
	}

	return user.ID, nil
}

func (u *UserUsecase) UpdateUser(id int, payload *domain.UserUpdateRequest) error {
	user, err := u.repo.GetUserBy("id", id)
	if err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserUpdateFailed, err)
	}

	if payload.Username != nil {
		existingUser, err := u.repo.CheckExistsBy("username", *payload.Username, &id)
		if err == nil && existingUser {
			return utils.ErrUsernameAlreadyExists
		}
		user.Username = *payload.Username
	}

	if payload.Email != nil {
		existingUser, err := u.repo.CheckExistsBy("email", *payload.Email, &id)
		if err == nil && existingUser {
			return utils.ErrEmailAlreadyExists
		}
		user.Email = *payload.Email
	}

	if err := u.repo.UpdateUser(id, &user); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserUpdateFailed, err)
	}

	return nil
}

func (u *UserUsecase) ChangePasswordUser(id int, payload *domain.UserChangePasswordRequest) error {
	user, err := u.repo.GetUserBy("id", id)
	if err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserChangePasswordFailed, err)
	}

	if !utils.VerifyPassword(user.Password, payload.OldPassword) {
		return fmt.Errorf("%w: %v", utils.ErrInvalidRequest, "Password lama tidak sesuai")
	}

	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		return utils.ErrHashPassword
	}

	if err := u.repo.ChangePasswordUser(id, string(hashedPassword)); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserChangePasswordFailed, err)
	}

	return nil
}

func (u *UserUsecase) ChangeStatusUser(id int, payload *domain.UserChangeStatusRequest) error {
	if err := u.repo.ChangeStatusUser(id, payload.Status); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserChangeStatusFailed, err)
	}

	return nil
}

func (u *UserUsecase) GetUserById(id int) (domain.UserDetailResponse, error) {
	user, err := u.repo.GetUserBy("id", id)
	if err != nil {
		return domain.UserDetailResponse{}, fmt.Errorf("%w: %v", utils.ErrGetDataFailed, err)
	}

	return user.ToUserDetailResponse(), nil
}

func (u *UserUsecase) Login(payload *domain.LoginRequest) (domain.LoginResponse, error) {
	user, err := u.repo.GetUserBy("email", payload.Email)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("%w: %v", utils.ErrLoginFailed, "Email atau password tidak sesuai")
	}

	if !utils.VerifyPassword(user.Password, payload.Password) {
		return domain.LoginResponse{}, fmt.Errorf("%w: %v", utils.ErrLoginFailed, "Email atau password tidak sesuai")
	}

	if user.Status != domain.UserActive {
		return domain.LoginResponse{}, fmt.Errorf("%w: %v", utils.ErrLoginFailed, "Akun tidak aktif")
	}

	token, err := token.GenerateToken(user.ID, user.Email)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("%w: %v", utils.ErrTokenCreatedFailed, err)
	}

	return domain.LoginResponse{Token: token}, nil
}
