package usecase

import (
	"database/sql"
	"fmt"

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

func (u *ProductUsecase) GetProductById(id int) (domain.ProductDetailResponse, error) {
	product, err := u.productRepo.GetProductById(id)
	if err != nil {
		return domain.ProductDetailResponse{}, fmt.Errorf("%w: %v", domain.ErrProductNotFound, err)
	}
	return product.ToProductResponse(), nil
}

func (u *ProductUsecase) GetAllProducts(filter domain.ProductFilter) ([]domain.ProductDetailResponse, error) {
	products, err := u.productRepo.GetAllProducts(filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrProductQueryFailed, err)
	}

	response := make([]domain.ProductDetailResponse, 0, len(products))
	for _, p := range products {
		response = append(response, p.ToProductResponse())
	}

	return response, nil
}

func (u *ProductUsecase) CreateProduct(payload *domain.ProductCreateRequest) error {
	product := domain.Product{
		Name:   payload.Name,
		Price:  payload.Price,
		Stock:  payload.Stock,
		Status: payload.Status,
	}

	if err := u.productRepo.CreateProduct(&product); err != nil {
		return domain.ErrProductCreationFailed
	}

	return nil

}

func (u *ProductUsecase) UpdateProduct(id int, payload *domain.ProductUpdateRequest) error {
	err := u.productRepo.UpdateProduct(id, payload)
	if err == sql.ErrNoRows {
		return domain.ErrProductNotFound
	} else if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrProductUpdateFailed, err)
	} else {
		return nil
	}
}

func (u *ProductUsecase) DeleteProduct(id int) error {
	err := u.productRepo.DeleteProduct(id)
	if err == sql.ErrNoRows {
		return domain.ErrProductNotFound
	} else if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrProductDeleteFailed, err)
	} else {
		return nil
	}
}
