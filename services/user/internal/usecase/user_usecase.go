package usecase

import (
	"database/sql"
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

func (u *UserUsecase) CreateUser(payload *domain.UserCreateRequest) error {
	existingUser, err := u.repo.CheckExistsBy("username", payload.Username)
	if err == nil && existingUser {
		return utils.ErrUsernameAlreadyExists
	}

	existingUser, err = u.repo.CheckExistsBy("email", payload.Email)
	if err == nil && existingUser {
		return utils.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return utils.ErrHashPassword
	}

	user := domain.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashedPassword),
		Status:   domain.UserActive,
	}

	if err := u.repo.CreateUser(&user); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserCreationFailed, err)
	}

	return nil
}

func (u *UserUsecase) UpdateUser(id int, payload *domain.UserUpdateRequest) error {
	user, err := u.repo.GetUserById(id)
	if err != nil {
		return utils.ErrUserNotFound
	}

	if payload.Username != nil {
		existingUser, err := u.repo.CheckExistsBy("username", *payload.Username)
		if err == nil && existingUser {
			return utils.ErrUsernameAlreadyExists
		}
		user.Username = *payload.Username
	}

	if payload.Email != nil {
		existingUser, err := u.repo.CheckExistsBy("email", *payload.Email)
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
	user, err := u.repo.GetUserById(id)
	if err != nil {
		return utils.ErrUserNotFound
	}

	if !utils.VerifyPassword(user.Password, payload.OldPassword) {
		return fmt.Errorf("%w: %v", utils.ErrInvalidRequest, "Password lama tidak sesuai")
	}

	if payload.NewPassword != payload.ConfirmPassword {
		return fmt.Errorf("%w: %v", utils.ErrInvalidRequest, "Konfirmasi password tidak sama")
	}

	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		return utils.ErrHashPassword
	}

	if err := u.repo.ChangePasswordUser(id, string(hashedPassword)); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserUpdateFailed, err)
	}

	return nil
}

func (u *UserUsecase) ChangeStatusUser(id int, payload *domain.UserChangeStatusRequest) error {
	err := u.repo.ChangeStatusUser(id, payload.Status)
	if err == sql.ErrNoRows {
		return utils.ErrUserNotFound
	} else if err != nil {
		return fmt.Errorf("%w: %v", utils.ErrUserUpdateFailed, err)
	} else {
		return nil
	}
}

func (u *UserUsecase) GetUserById(id int) (domain.UserDetailResponse, error) {
	user, err := u.repo.GetUserById(id)
	if err != nil {
		return domain.UserDetailResponse{}, fmt.Errorf("%w: %v", utils.ErrUserNotFound, err)
	}

	return user.ToUserDetailResponse(), nil
}

func (u *UserUsecase) Login(payload *domain.LoginRequest) (domain.LoginResponse, error) {
	user, err := u.repo.GetUserByEmail(payload.Email)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	if !utils.VerifyPassword(user.Password, payload.Password) {
		return domain.LoginResponse{}, utils.ErrInvalidRequest
	}

	token, err := token.GenerateToken(user.ID, user.Email)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("%w: %v", utils.ErrTokenCreatedFailed, err)
	}

	return domain.LoginResponse{Token: token}, nil
}
