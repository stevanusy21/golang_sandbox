package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/stevanusy21/golang_sandbox/pkg/request"
	"github.com/stevanusy21/golang_sandbox/pkg/response"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/usecase"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productUsecase.GetProductById(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	allowedColumns := map[string]bool{
		"id":     true,
		"name":   true,
		"price":  true,
		"stock":  true,
		"status": true,
	}

	statusQuery := r.URL.Query().Get("status")
	statusSlice := []string{}

	if statusQuery != "" {
		statusSlice = strings.Split(statusQuery, ",")
	}

	filter := domain.ProductFilter{
		Id:     r.URL.Query().Get("id"),
		Name:   r.URL.Query().Get("name"),
		Price:  r.URL.Query().Get("price"),
		Status: statusSlice,
		Pagination: utils.GeneratePaginationData(
			r.URL.Query().Get("page"),
			r.URL.Query().Get("limit"),
			r.URL.Query().Get("sort_by"),
			r.URL.Query().Get("sort_dir"),
			allowedColumns,
		),
	}

	products, err := h.productUsecase.GetAllProducts(filter)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductQueryFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	payload, err := request.DecodeJSON[domain.ProductCreateRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.productUsecase.CreateProduct(&payload); err != nil {
		switch {
		case errors.Is(err, domain.ErrProductCreationFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "Produk berhasil dibuat"})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, err := request.DecodeJSON[domain.ProductUpdateRequest](r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.productUsecase.UpdateProduct(id, &payload); err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrProductUpdateFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Produk berhasil diupdate"})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.productUsecase.DeleteProduct(id); err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrProductDeleteFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal pada server")
		}
		return
	}

	response.JSON(w, http.StatusOK, "Produk berhasil di hapus!")
}
