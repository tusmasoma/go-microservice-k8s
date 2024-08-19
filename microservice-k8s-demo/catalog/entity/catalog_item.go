package entity

import (
	"errors"

	"github.com/google/uuid"
)

type CatalogItem struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewCatalogItem(name string, price float64) (*CatalogItem, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	return &CatalogItem{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}, nil
}
