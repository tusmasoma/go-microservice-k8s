//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/entity"
)

type CatalogItemRepository interface {
	Get(ctx context.Context, id string) (*entity.CatalogItem, error)
	List(ctx context.Context) ([]entity.CatalogItem, error)
	Create(ctx context.Context, task entity.CatalogItem) error
	Update(ctx context.Context, task entity.CatalogItem) error
	Delete(ctx context.Context, id string) error
}
