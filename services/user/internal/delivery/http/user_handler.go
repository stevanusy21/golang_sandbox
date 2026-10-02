package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/stevanusy21/golang_sandbox/pkg/request"
	"github.com/stevanusy21/golang_sandbox/pkg/response"
	"github.com/stevanusy21/golang_sandbox/pkg/router"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	r := router.NewRouter(mux)

	r.Public("POST /auth/login", h.Login)
	r.Public("POST /users", h.CreateUser)
	r.Public("PUT /users/{id}/password", h.ChangePasswordUser)
	r.Protected("GET /users/{id}", h.GetUserById)
	r.Protected("PATCH /users/{id}/profile", h.UpdateUser)
	r.Protected("PATCH /users/{id}/status", h.ChangeStatusUser)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	payload, err := request.DecodeJSON[domain.UserCreateRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.userUsecase.CreateUser(&payload)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrUsernameAlreadyExists), errors.Is(err, utils.ErrEmailAlreadyExists):
			response.Error(w, http.StatusConflict, err.Error())
		case errors.Is(err, utils.ErrHashPassword):
			response.Error(w, http.StatusInternalServerError, err.Error())
		case errors.Is(err, utils.ErrUserCreationFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "User berhasil dibuat",
		"user_id": strconv.Itoa(userId),
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, err := request.DecodeJSON[domain.UserUpdateRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userUsecase.UpdateUser(id, &payload); err != nil {
		switch {
		case errors.Is(err, utils.ErrUsernameAlreadyExists), errors.Is(err, utils.ErrEmailAlreadyExists):
			response.Error(w, http.StatusConflict, err.Error())
		case errors.Is(err, utils.ErrUserUpdateFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "User berhasil diupdate"})
}

func (h *UserHandler) ChangePasswordUser(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, err := request.DecodeJSON[domain.UserChangePasswordRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userUsecase.ChangePasswordUser(id, &payload); err != nil {
		switch {
		case errors.Is(err, utils.ErrHashPassword):
			response.Error(w, http.StatusInternalServerError, err.Error())
		case errors.Is(err, utils.ErrInvalidRequest):
			response.Error(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, utils.ErrUserChangePasswordFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Password berhasil diubah"})
}

func (h *UserHandler) ChangeStatusUser(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, err := request.DecodeJSON[domain.UserChangeStatusRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userUsecase.ChangeStatusUser(id, &payload); err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, utils.ErrUserChangeStatusFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Status berhasil diubah"})
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userUsecase.GetUserById(id)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrUserGetFailed):
			response.Error(w, http.StatusNotFound, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	payload, err := request.DecodeJSON[domain.LoginRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.userUsecase.Login(&payload)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrLoginFailed):
			response.Error(w, http.StatusUnauthorized, err.Error())
		case errors.Is(err, utils.ErrTokenCreatedFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, result)
}
