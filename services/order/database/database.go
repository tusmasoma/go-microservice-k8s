//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package database

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/services/order/entity"
)

type Database struct {
	Order Order
}

type Order interface {
	Get(ctx context.Context, id string) (*entity.Order, error)
	List(ctx context.Context) (entity.Orders, error)
	Create(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id string) error
}
