package usecase

import (
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
		return domain.ProductDetailResponse{}, err
	}
	return product.ToProductResponse(), nil
}

func (u *ProductUsecase) GetAllProducts(filter domain.ProductFilter) (utils.PaginationResponse, error) {
	products, total, err := u.productRepo.GetAllProducts(filter)
	if err != nil {
		return utils.PaginationResponse{}, err
	}

	response := make([]domain.ProductDetailResponse, 0, len(products))
	for _, p := range products {
		response = append(response, p.ToProductResponse())
	}

	return utils.NewPaginationResponse(response, total, filter.Page, filter.Limit), nil
}

func (u *ProductUsecase) CreateProduct(payload *domain.ProductCreateRequest) (int, error) {
	product := domain.Product{
		Name:   payload.Name,
		Price:  payload.Price,
		Stock:  payload.Stock,
		Status: payload.Status,
	}

	if err := u.productRepo.CreateProduct(&product); err != nil {
		return 0, err
	}

	return product.ID, nil

}

func (u *ProductUsecase) UpdateProduct(id int, payload *domain.ProductUpdateRequest) error {
	return u.productRepo.UpdateProduct(id, payload)
}

func (u *ProductUsecase) DeleteProduct(id int) error {
	return u.productRepo.DeleteProduct(id)
}

func (u *ProductUsecase) DeductStock(id int, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("%w: %v", utils.ErrInvalidRequest, "Quantity harus lebih besar dari 0")
	}

	return u.productRepo.DeductStock(id, quantity)
}
