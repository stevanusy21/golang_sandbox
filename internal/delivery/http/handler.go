package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stevanusy21/golang_sandbox/internal/domain"
	"github.com/stevanusy21/golang_sandbox/internal/usecase"
	"github.com/stevanusy21/golang_sandbox/pkg/response"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type Handler struct {
	productUsecase *usecase.ProductUsecase
}

func NewHandler(productUsecase *usecase.ProductUsecase) *Handler {
	return &Handler{
		productUsecase: productUsecase,
	}
}

func (h *Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "ID produk harus diisi")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID produk tidak valid")
		return
	}

	product, err := h.productUsecase.GetProductById(id)
	if err == sql.ErrNoRows {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Data product dengan id %d tidak ditemukan", id))
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	allowedColumns := map[string]bool{
		"id": true,
		"name": true,
		"price": true,
	}
	
	filter := domain.ProductFilter{
		ID: r.URL.Query().Get("id"),
		Name: r.URL.Query().Get("name"),
		Price: r.URL.Query().Get("price"),
		Pagination: utils.GeneratePagination(
			r.URL.Query().Get("page"),
			r.URL.Query().Get("limit"),
			r.URL.Query().Get("sort_by"),
			r.URL.Query().Get("sort_dir"),
			allowedColumns,
		),
	}
	
	products, err := h.productUsecase.GetAllProducts(filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, products)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input domain.Product

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	err = h.productUsecase.CreateProduct(&input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, input)
}