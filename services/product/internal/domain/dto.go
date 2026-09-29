package domain

import (
	"errors"
	"time"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type ProductFilter struct {
	Id     string
	Name   string
	Price  string
	Status []string
	utils.Pagination
}

type ProductCreateRequest struct {
	Name   string        `json:"name"`
	Price  float64       `json:"price"`
	Stock  int           `json:"stock"`
	Status ProductStatus `json:"status"`
}

type ProductUpdateRequest struct {
	Name   *string        `json:"name"`
	Price  *float64       `json:"price"`
	Stock  *int           `json:"stock"`
	Status *ProductStatus `json:"status"`
}

type ProductDetailResponse struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Price     float64       `json:"price"`
	Stock     int           `json:"stock"`
	Status    ProductStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

func (r ProductCreateRequest) Validate() error {
	if utils.IsEmpty(r.Name) {
		return errors.New("Name tidak boleh kosong")
	}

	if utils.IsEmpty(r.Price) {
		return errors.New("Price tidak boleh kosong")
	}

	if utils.IsEmpty(r.Stock) {
		return errors.New("Stock tidak boleh kosong")
	}

	if utils.IsEmpty(r.Status) {
		return errors.New("Status tidak boleh kosong")
	}

	if !r.Status.IsValid() {
		return errors.New("Status tidak valid")
	}

	if r.Price < 0 {
		return errors.New("Harga tidak boleh kurang dari 0")
	}

	if r.Stock < 0 {
		return errors.New("Stock tidak boleh kurang dari 0")
	}

	return nil
}

func (r ProductUpdateRequest) Validate() error {
	if r.Name == nil && r.Price == nil && r.Stock == nil && r.Status == nil {
		return errors.New("Tidak ada data yang diubah")
	}

	if r.Name != nil && utils.IsEmpty(*r.Name) {
		return errors.New("Name tidak boleh kosong")
	}

	if r.Price != nil && *r.Price < 0 {
		return errors.New("Harga tidak boleh kurang dari 0")
	}

	if r.Stock != nil && *r.Stock < 0 {
		return errors.New("Stock tidak boleh kurang dari 0")
	}

	if r.Status != nil && !(*r.Status).IsValid() {
		return errors.New("Status tidak valid")
	}

	return nil
}
