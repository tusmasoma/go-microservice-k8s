package entity

import (
	"errors"

	"github.com/google/uuid"
)

type CatalogItem struct {
	ID    string  `json:"id" db:"id"`
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}

func NewCatalogItem(id, name string, price float64) (*CatalogItem, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	return &CatalogItem{
		ID:    id,
		Name:  name,
		Price: price,
	}, nil
}
