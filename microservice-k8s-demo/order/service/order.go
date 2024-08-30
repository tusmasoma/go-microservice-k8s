//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package service

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order entity.Order) error
	DeleteOrder(ctx context.Context, id string) error
}

type orderService struct {
	or  repository.OrderRepository
	olr repository.OrderLineRepository
	tr  repository.TransactionRepository
}

func NewOrderService(
	or repository.OrderRepository,
	olr repository.OrderLineRepository,
	tr repository.TransactionRepository,
) OrderService {
	return &orderService{
		or:  or,
		olr: olr,
		tr:  tr,
	}
}

func (os *orderService) CreateOrder(ctx context.Context, order entity.Order) error {
	if err := os.tr.Transaction(ctx, func(ctx context.Context) error {
		if err := os.or.Create(ctx, entity.OrderModel{
			ID:         order.ID,
			CustomerID: order.Customer.ID,
			OrderDate:  order.OrderDate,
		}); err != nil {
			return err
		}

		var orderLiensModel []entity.OrderLineModel
		for _, ol := range order.OrderLines {
			orderLiensModel = append(orderLiensModel, entity.OrderLineModel{
				OrderID:       order.ID,
				CatalogItemID: ol.CatalogItem.ID,
				Count:         ol.Count,
			})
		}

		if err := os.olr.BatchCreate(ctx, orderLiensModel); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (os *orderService) DeleteOrder(ctx context.Context, id string) error {
	if err := os.tr.Transaction(ctx, func(ctx context.Context) error {
		if err := os.olr.BatchDelete(ctx, id); err != nil {
			return err
		}
		if err := os.or.Delete(ctx, id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
