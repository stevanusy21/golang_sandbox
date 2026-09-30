package grpc

import (
	"context"

	"github.com/stevanusy21/golang_sandbox/proto/product"
	"github.com/stevanusy21/golang_sandbox/services/product/internal/usecase"
)

type ProductHandlerGrpc struct {
	product.UnimplementedProductServiceServer
	productUseCase *usecase.ProductUsecase
}

func NewProductHandlerGrpc(usecase *usecase.ProductUsecase) *ProductHandlerGrpc {
	return &ProductHandlerGrpc{productUseCase: usecase}
}

func (h *ProductHandlerGrpc) GetProductDetail(ctx context.Context, req *product.ProductDetailRequest) (*product.ProductDetailResponse, error) {
	p, err := h.productUseCase.GetProductById(int(req.ProductId))
	if err != nil {
		return nil, err
	}
	return &product.ProductDetailResponse{
		ProductId: int32(p.ID),
		Name:      p.Name,
		Price:     float64(p.Price),
		Stock:     int32(p.Stock),
	}, nil
}

func (h *ProductHandlerGrpc) DeductStock(ctx context.Context, req *product.DeductStockRequest) (*product.DeductStockResponse, error) {
	err := h.productUseCase.DeductStock(int(req.ProductId), int(req.Quantity))
	if err != nil {
		return &product.DeductStockResponse{
			IsSuccess: false,
			Message:   err.Error(),
		}, nil
	}
	return &product.DeductStockResponse{
		IsSuccess: true,
		Message:   "Stock berhasil dikurangi",
	}, nil
}
