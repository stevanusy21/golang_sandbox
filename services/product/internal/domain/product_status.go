package domain

type ProductStatus string

const (
	Available    ProductStatus = "AVAILABLE"
	OutOfStock   ProductStatus = "OUT_OF_STOCK"
	Discontinued ProductStatus = "DISCONTINUED"
)

func (s ProductStatus) IsValid() bool {
	switch s {
	case Available, OutOfStock, Discontinued:
		return true
	default:
		return false
	}
}
