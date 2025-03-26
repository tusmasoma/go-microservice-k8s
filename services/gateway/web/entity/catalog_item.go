package entity

import (
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	"github.com/tusmasoma/go-microservice-k8s/proto/web"
)

type CatalogItem struct {
	web.CatalogItem
}

func (c *CatalogItem) Proto() *web.CatalogItem {
	if c == nil {
		return nil
	}
	return &c.CatalogItem
}

func NewCatalogItem(item *catalog.CatalogItem) *CatalogItem {
	return &CatalogItem{
		CatalogItem: web.CatalogItem{
			Id:    item.GetId(),
			Name:  item.GetName(),
			Price: item.GetPrice(),
		},
	}
}

type CatalogItems []*CatalogItem

func NewCatalogItems(items []*catalog.CatalogItem) CatalogItems {
	ret := make(CatalogItems, len(items))
	for i := range items {
		ret[i] = NewCatalogItem(items[i])
	}
	return ret
}

func (cs CatalogItems) Proto() []*web.CatalogItem {
	items := make([]*web.CatalogItem, len(cs))
	for i, item := range cs {
		items[i] = item.Proto()
	}
	return items
}
