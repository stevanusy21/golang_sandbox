package domain

import "github.com/stevanusy21/golang_sandbox/pkg/utils"

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	Status string `json:"status"`
}

type ProductFilter struct {
	ID    string
	Name  string
	Price string
	Status []string
	utils.Pagination
}