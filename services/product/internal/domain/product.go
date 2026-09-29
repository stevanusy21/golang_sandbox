package domain

import "time"

type Product struct {
	ID        int           `db:"id"`
	Name      string        `db:"name"`
	Price     float64       `db:"price"`
	Stock     int           `db:"stock"`
	Status    ProductStatus `db:"status"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	DeletedAt *time.Time    `db:"deleted_at"`
}

func (p Product) ToProductResponse() ProductDetailResponse {
	return ProductDetailResponse{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
	}
}
