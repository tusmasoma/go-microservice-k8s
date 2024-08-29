//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
)

type OrderRepository interface {
	Get(ctx context.Context, id string) (*entity.OrderModel, error)
	List(ctx context.Context) ([]entity.OrderModel, error)
	Create(ctx context.Context, order entity.OrderModel) error
	Delete(ctx context.Context, id string) error
}

type OrderLineRepository interface {
	List(ctx context.Context, orderID string) ([]entity.OrderLineModel, error)
	Create(ctx context.Context, orderLine entity.OrderLineModel) error
	BatchCreate(ctx context.Context, orderLines []entity.OrderLineModel) error
	Delete(ctx context.Context, orderID, catalogItemID string) error
	BatchDelete(ctx context.Context, orderID string) error
}
