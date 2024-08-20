//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/entity"
	"github.com/tusmasoma/microservice-k8s-demo/catalog/repository"
)

type CatalogItemUseCase interface {
	GetCatalogItem(ctx context.Context, id string) (*entity.CatalogItem, error)
	ListCatalogItems(ctx context.Context) ([]entity.CatalogItem, error)
	ListCatalogItemsByName(ctx context.Context, name string) ([]entity.CatalogItem, error)
	CreateCatalogItem(ctx context.Context, name string, price float64) error
	UpdateCatalogItem(ctx context.Context, id, name string, price float64) error
	DeleteCatalogItem(ctx context.Context, id string) error
}

type catalogItemUseCase struct {
	cr repository.CatalogItemRepository
}

func NewCatalogItemUseCase(cr repository.CatalogItemRepository) CatalogItemUseCase {
	return &catalogItemUseCase{
		cr: cr,
	}
}

func (cu *catalogItemUseCase) GetCatalogItem(ctx context.Context, id string) (*entity.CatalogItem, error) {
	item, err := cu.cr.Get(ctx, id)
	if err != nil {
		log.Error("Failed to get catalog item", log.Ferror(err))
		return nil, err
	}
	return item, nil
}

func (cu *catalogItemUseCase) ListCatalogItems(ctx context.Context) ([]entity.CatalogItem, error) {
	items, err := cu.cr.List(ctx)
	if err != nil {
		log.Error("Failed to list catalog items", log.Ferror(err))
		return nil, err
	}
	return items, nil
}

func (cu *catalogItemUseCase) ListCatalogItemsByName(ctx context.Context, name string) ([]entity.CatalogItem, error) {
	items, err := cu.cr.ListByName(ctx, name)
	if err != nil {
		log.Error("Failed to list catalog items by name", log.Ferror(err))
		return nil, err
	}
	return items, nil
}

func (cu *catalogItemUseCase) CreateCatalogItem(ctx context.Context, name string, price float64) error {
	item, err := entity.NewCatalogItem(name, price)
	if err != nil {
		log.Error("Failed to create catalog item", log.Ferror(err))
		return err
	}
	if err = cu.cr.Create(ctx, *item); err != nil {
		log.Error("Failed to create catalog item", log.Ferror(err))
		return err
	}
	return nil
}

func (cu *catalogItemUseCase) UpdateCatalogItem(ctx context.Context, id, name string, price float64) error {
	item, err := cu.cr.Get(ctx, id)
	if err != nil {
		log.Error("Failed to get catalog item", log.Ferror(err))
		return err
	}

	item.Name = name
	item.Price = price

	if err = cu.cr.Update(ctx, *item); err != nil {
		log.Error("Failed to update catalog item", log.Ferror(err))
		return err
	}
	return nil
}

func (cu *catalogItemUseCase) DeleteCatalogItem(ctx context.Context, id string) error {
	if err := cu.cr.Delete(ctx, id); err != nil {
		log.Error("Failed to delete catalog item", log.Ferror(err))
		return err
	}
	return nil
}
