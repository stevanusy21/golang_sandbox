package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/stevanusy21/golang_sandbox/pkg/request"
	"github.com/stevanusy21/golang_sandbox/pkg/response"
	"github.com/stevanusy21/golang_sandbox/pkg/router"
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

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	r := router.NewRouter(mux)

	r.Protected("GET /products", h.GetAllProducts)
	r.Protected("GET /products/{id}", h.GetProductByID)
	r.Protected("POST /products", h.CreateProduct)
	r.Protected("PUT /products/{id}", h.UpdateProduct)
	r.Protected("DELETE /products/{id}", h.DeleteProduct)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id", true)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productUsecase.GetProductById(id)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, utils.ErrQueryFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, utils.ErrInternalServerError.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	allowedSortColumns := map[string]bool{
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
			allowedSortColumns,
		),
	}

	products, err := h.productUsecase.GetAllProducts(filter)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrQueryFailed), errors.Is(err, utils.ErrScanFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, utils.ErrInternalServerError.Error())
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

	productId, err := h.productUsecase.CreateProduct(&payload)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrQueryFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, utils.ErrInternalServerError.Error())
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"message":   "Produk berhasil dibuat",
		"productId": strconv.Itoa(productId),
	})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id", true)
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
		case errors.Is(err, utils.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, utils.ErrQueryFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, utils.ErrInternalServerError.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Produk berhasil diupdate"})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIntParam(r, "id", true)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.productUsecase.DeleteProduct(id); err != nil {
		switch {
		case errors.Is(err, utils.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, utils.ErrQueryFailed):
			response.Error(w, http.StatusInternalServerError, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, utils.ErrInternalServerError.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, "Produk berhasil di hapus!")
}
