package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
)

type CatalogItem struct {
	catalog.CatalogItem
}

func (c *CatalogItem) Proto() *catalog.CatalogItem {
	if c == nil {
		return nil
	}
	return &c.CatalogItem
}

func NewCatalogItem(id, name string, price float64) (*CatalogItem, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	return &CatalogItem{
		CatalogItem: catalog.CatalogItem{
			Id:    id,
			Name:  name,
			Price: price,
		},
	}, nil
}

func CreateCatalogItem(name string, price float64) (*CatalogItem, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	return &CatalogItem{
		CatalogItem: catalog.CatalogItem{
			Id:    uuid.NewString(),
			Name:  name,
			Price: price,
		},
	}, nil
}

type CatalogItems []*CatalogItem

func (cs CatalogItems) Proto() []*catalog.CatalogItem {
	items := make([]*catalog.CatalogItem, len(cs))
	for i, item := range cs {
		items[i] = item.Proto()
	}
	return items
}
