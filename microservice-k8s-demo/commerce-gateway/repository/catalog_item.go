//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/commerce-gateway/entity"
)

type CatalogItemRepository interface {
	List(ctx context.Context) ([]entity.CatalogItem, error)
	ListByName(ctx context.Context, name string) ([]entity.CatalogItem, error)
	Create(ctx context.Context, item entity.CatalogItem) error
	Update(ctx context.Context, item entity.CatalogItem) error
	Delete(ctx context.Context, id string) error
}
