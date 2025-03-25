//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package database

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/services/catalog/entity"
)

type Database struct {
	CatalogItem CatalogItem
}

type CatalogItem interface {
	Get(ctx context.Context, id string) (*entity.CatalogItem, error)
	List(ctx context.Context) (entity.CatalogItems, error)
	ListByName(ctx context.Context, name string) (entity.CatalogItems, error)
	ListByIDs(ctx context.Context, ids []string) (entity.CatalogItems, error)
	Create(ctx context.Context, item *entity.CatalogItem) error
	Update(ctx context.Context, item *entity.CatalogItem) error
	Delete(ctx context.Context, id string) error
}
