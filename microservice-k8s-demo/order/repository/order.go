//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
)

type OrderRepository interface {
	Get(ctx context.Context, id string) (*entity.Order, error)
	List(ctx context.Context) ([]*entity.Order, error)
	Create(ctx context.Context, order entity.Order) error
	Delete(ctx context.Context, id string) error
}
