//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type OrderUseCase interface {
	GetOrderCreationResources(ctx context.Context) ([]entity.Customer, []entity.CatalogItem, error)
	GetOrder(ctx context.Context, id string) (*entity.Order, error)
	ListOrders(ctx context.Context) ([]*entity.Order, error)
	CreateOrder(ctx context.Context, params *CreateOrderParams) error
	DeleteOrder(ctx context.Context, id string) error
}

type orderUseCase struct {
	cr  repository.CustomerRepository
	cir repository.CatalogItemRepository
	or  repository.OrderRepository
}

func NewOrderUseCase(
	cr repository.CustomerRepository,
	cir repository.CatalogItemRepository,
	or repository.OrderRepository,
) OrderUseCase {
	return &orderUseCase{
		cr:  cr,
		cir: cir,
		or:  or,
	}
}

func (ouc *orderUseCase) GetOrderCreationResources(ctx context.Context) ([]entity.Customer, []entity.CatalogItem, error) {
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

func (ouc *orderUseCase) GetOrder(ctx context.Context, id string) (*entity.Order, error) {
	order, err := ouc.or.Get(ctx, id)
	if err != nil {
		log.Error("Failed to get order", log.Ferror(err))
		return nil, err
	}

	customer, err := ouc.cr.Get(ctx, order.Customer.ID)
	if err != nil {
		log.Error("Failed to get customer", log.Ferror(err))
		return nil, err
	}

	var orderLines []entity.OrderLine
	for _, ol := range order.OrderLines {
		// TODO: N + 1 problem
		item, err := ouc.cir.Get(ctx, ol.CatalogItem.ID)
		if err != nil {
			return nil, err
		}
		orderLines = append(orderLines, entity.OrderLine{
			Count:       ol.Count,
			CatalogItem: *item,
		})
	}

	order.Customer = *customer
	order.OrderLines = orderLines

	order.TotalPrice = order.GetTotalPrice()

	return order, nil
}

func (ouc *orderUseCase) ListOrders(ctx context.Context) ([]*entity.Order, error) {
	orders, err := ouc.or.List(ctx)
	if err != nil {
		log.Error("Failed to get orders", log.Ferror(err))
		return nil, err
	}

	for _, order := range orders {
		customer, err := ouc.cr.Get(ctx, order.Customer.ID)
		if err != nil {
			log.Error("Failed to get customer", log.Ferror(err))
			return nil, err
		}

		var orderLines []entity.OrderLine
		for _, ol := range order.OrderLines {
			// TODO: N + 1 problem
			item, err := ouc.cir.Get(ctx, ol.CatalogItem.ID)
			if err != nil {
				return nil, err
			}
			orderLines = append(orderLines, entity.OrderLine{
				Count:       ol.Count,
				CatalogItem: *item,
			})
		}

		order.Customer = *customer
		order.OrderLines = orderLines
		order.TotalPrice = order.GetTotalPrice()
	}

	return orders, nil
}

type CreateOrderParams struct {
	CustomerID string
	OrderLine  []struct {
		CatalogItemID string
		Count         int
	}
}

func (ouc *orderUseCase) CreateOrder(ctx context.Context, params *CreateOrderParams) error {
	customer, err := ouc.cr.Get(ctx, params.CustomerID)
	if err != nil {
		log.Error("Failed to get cusotmer", log.Ferror(err))
		return err
	}

	var orderLiens []entity.OrderLine
	for _, ol := range params.OrderLine {
		item, err := ouc.cir.Get(ctx, ol.CatalogItemID) //nolint:govet // err shadow
		if err != nil {
			log.Error("Failed to get catalog item", log.Ferror(err))
			return err
		}
		orderLine, err := entity.NewOrderLine(ol.Count, *item)
		if err != nil {
			log.Error("Failed to create order line", log.Ferror(err))
			return err
		}
		orderLiens = append(orderLiens, *orderLine)
	}

	order, err := entity.NewOrder(*customer, orderLiens)
	if err != nil {
		log.Error("Failed to create order", log.Ferror(err))
		return err
	}
	if err = ouc.or.Create(ctx, *order); err != nil {
		log.Error("Failed to create order", log.Ferror(err))
		return err
	}
	return nil
}

func (ouc *orderUseCase) DeleteOrder(ctx context.Context, id string) error {
	if err := ouc.or.Delete(ctx, id); err != nil {
		log.Error("Failed to delete order", log.Ferror(err))
		return err
	}
	return nil
}
