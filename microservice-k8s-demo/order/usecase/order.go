//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type OrderUseCase interface {
	GetOrderPageData(ctx context.Context) ([]entity.Customer, []entity.CatalogItem, error)
}

type orderUseCase struct {
	cr  repository.CustomerRepository
	cir repository.CatalogItemRepository
	or  repository.OrderRepository
}

func NewOrderUseCase(cr repository.CustomerRepository, cir repository.CatalogItemRepository, or repository.OrderRepository) OrderUseCase {
	return &orderUseCase{
		cr:  cr,
		cir: cir,
		or:  or,
	}
}

func (ouc *orderUseCase) GetOrderPageData(ctx context.Context) ([]entity.Customer, []entity.CatalogItem, error) {
	customers, err := ouc.cr.List(ctx)
	if err != nil {
		log.Error("Failed to get customer", log.Ferror(err))
		return nil, nil, err
	}
	items, err := ouc.cir.List(ctx)
	if err != nil {
		log.Error("Failed to get catalog item", log.Ferror(err))
		return nil, nil, err
	}
	return customers, items, nil
}