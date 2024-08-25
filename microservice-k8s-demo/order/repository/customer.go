//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
)

type CustomerRepository interface {
	List(ctx context.Context) ([]entity.Customer, error)
	Create(ctx context.Context, customer entity.Customer) error
	Update(ctx context.Context, customer entity.Customer) error
	Delete(ctx context.Context, id string) error
}
