package usecase

import (
	"errors"
	"strings"

	"github.com/stevanusy21/golang_sandbox/internal/domain"
	"github.com/stevanusy21/golang_sandbox/internal/repository"
)

type ProductUsecase struct {
	productRepo *repository.ProductRepository
}

func NewProductUsecase(repo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: repo,
	}
}

func (u *ProductUsecase) GetAllProducts(filter domain.ProductFilter) ([]domain.Product, error) {
	return u.productRepo.GetAllProducts(filter)
}

func (u *ProductUsecase) CreateProduct(product *domain.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("Nama produk tidak boleh kosong")
	}

	if product.Price <= 0 {
		return errors.New("Harga produk minimal 0 rupiah")
	} 

	return u.productRepo.CreateProduct(product)
}