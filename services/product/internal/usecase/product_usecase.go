package usecase

import (
	"errors"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
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
		if errors.Is(err, utils.ErrProductNotFound) {
			return domain.ProductDetailResponse{}, utils.ErrProductNotFound
		}
		return domain.ProductDetailResponse{}, fmt.Errorf("%w: %v", utils.ErrGetDataFailed, err)
	}
	return product.ToProductResponse(), nil
}

func (u *ProductUsecase) GetAllProducts(filter domain.ProductFilter) ([]domain.ProductDetailResponse, error) {
	products, err := u.productRepo.GetAllProducts(filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrGetDataFailed, err)
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
		return fmt.Errorf("%w: %v", utils.ErrProductCreationFailed, err)
	}

	return nil

}

func (u *ProductUsecase) UpdateProduct(id int, payload *domain.ProductUpdateRequest) error {
	if err := u.productRepo.UpdateProduct(id, payload); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrProductUpdateFailed, err)
	}
	return nil
}

func (u *ProductUsecase) DeleteProduct(id int) error {
	if err := u.productRepo.DeleteProduct(id); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrProductDeleteFailed, err)
	}
	return nil
}

func (u *ProductUsecase) DeductStock(id int, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("%w: %v", utils.ErrInvalidRequest, "Quantity harus lebih besar dari 0")
	}
	if err := u.productRepo.DeductStock(id, quantity); err != nil {
		return fmt.Errorf("%w: %v", utils.ErrDeductStockFailed, err)
	}

	return nil
}
