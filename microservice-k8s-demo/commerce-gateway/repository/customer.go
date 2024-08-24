//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/microservice-k8s-demo/commerce-gateway/entity"
)

type CustomerRepository interface {
	List(ctx context.Context) ([]entity.Customer, error)
	Create(ctx context.Context, customer entity.Customer) error
	Update(ctx context.Context, customer entity.Customer) error
	Delete(ctx context.Context, id string) error
}
