package usecase

import (
	"errors"
	"strings"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/domain"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/repository"
)

type ProductUsecase struct {
	productRepo *repository.ProductRepository
}

func NewProductUsecase(repo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: repo,
	}
}

func (u *ProductUsecase) GetProductById(id int) (domain.Product, error) {
	return u.productRepo.GetProductById(id)
}

func (u *ProductUsecase) GetAllProducts(filter domain.ProductFilter) ([]domain.Product, error) {
	return u.productRepo.GetAllProducts(filter)
}

func (u *ProductUsecase) CreateProduct(product *domain.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("Nama produk tidak boleh kosong")
	}
	if product.Price < 0 {
		return errors.New("Harga produk minimal 0 rupiah")
	}
	if product.Stock < 0 {
		return errors.New("Stock produk minimal 0")
	}

	return u.productRepo.CreateProduct(product)

}

func (u *ProductUsecase) UpdateProduct(product *domain.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("Nama produk tidak boleh kosong")
	}
	if product.Price < 0 {
		return errors.New("Harga produk minimal 0 rupiah")
	}
	if product.Stock < 0 {
		return errors.New("Stock produk minimal 0")
	}

	return u.productRepo.UpdateProduct(product)
}

func (u *ProductUsecase) DeleteProduct(id int) error {
	return u.productRepo.DeleteProduct(id)
}