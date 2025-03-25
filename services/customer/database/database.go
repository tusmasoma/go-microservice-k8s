//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package database

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/services/customer/entity"
)

type Database struct {
	Customer Customer
}

type Customer interface {
	Get(ctx context.Context, id string) (*entity.Customer, error)
	List(ctx context.Context) (entity.Customers, error)
	Create(ctx context.Context, customer *entity.Customer) error
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id string) error
}
