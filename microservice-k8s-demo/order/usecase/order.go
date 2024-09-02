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
	GetOrder(ctx context.Context, id string) (*OrderDetails, error)
	ListOrders(ctx context.Context) ([]*OrderDetails, error)
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

type OrderDetails struct {
	Order      *entity.Order
	Customer   *entity.Customer
	OrderLines []*OrderLineDetails
}

type OrderLineDetails struct {
	Count       int
	CatalogItem *entity.CatalogItem
}

func (ouc *orderUseCase) GetOrder(ctx context.Context, id string) (*OrderDetails, error) {
	order, err := ouc.or.Get(ctx, id)
	if err != nil {
		log.Error("Failed to get order", log.Ferror(err))
		return nil, err
	}

	customer, err := ouc.cr.Get(ctx, order.CustomerID)
	if err != nil {
		log.Error("Failed to get customer", log.Ferror(err))
		return nil, err
	}

	var orderLineDetails []*OrderLineDetails
	for _, ol := range order.OrderLines {
		// TODO: N + 1 problem
		item, err := ouc.cir.Get(ctx, ol.CatalogItemID)
		if err != nil {
			return nil, err
		}
		orderLineDetails = append(orderLineDetails, &OrderLineDetails{
			Count:       ol.Count,
			CatalogItem: item,
		})
	}

	return &OrderDetails{
		Order:      order,
		Customer:   customer,
		OrderLines: orderLineDetails,
	}, nil
}

func (ouc *orderUseCase) ListOrders(ctx context.Context) ([]*OrderDetails, error) {
	var orderDetails []*OrderDetails

	orders, err := ouc.or.List(ctx)
	if err != nil {
		log.Error("Failed to get orders", log.Ferror(err))
		return nil, err
	}

	for _, order := range orders {
		customer, err := ouc.cr.Get(ctx, order.CustomerID)
		if err != nil {
			log.Error("Failed to get customer", log.Ferror(err))
			return nil, err
		}

		var orderLineDetails []*OrderLineDetails
		for _, ol := range order.OrderLines {
			// TODO: N + 1 problem
			item, err := ouc.cir.Get(ctx, ol.CatalogItemID)
			if err != nil {
				return nil, err
			}
			orderLineDetails = append(orderLineDetails, &OrderLineDetails{
				Count:       ol.Count,
				CatalogItem: item,
			})
		}

		orderDetails = append(orderDetails, &OrderDetails{
			Order:      order,
			Customer:   customer,
			OrderLines: orderLineDetails,
		})
	}

	return orderDetails, nil
}

type CreateOrderParams struct {
	CustomerID string
	OrderLine  []struct {
		CatalogItemID string
		Count         int
	}
}

func (ouc *orderUseCase) CreateOrder(ctx context.Context, params *CreateOrderParams) error {
	var orderLiens []*entity.OrderLine
	for _, ol := range params.OrderLine {
		orderLine, err := entity.NewOrderLine(ol.Count, ol.CatalogItemID)
		if err != nil {
			log.Error("Failed to create order line", log.Ferror(err))
			return err
		}
		orderLiens = append(orderLiens, orderLine)
	}

	order, err := entity.NewOrder(params.CustomerID, orderLiens)
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
