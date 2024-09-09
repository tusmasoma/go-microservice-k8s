//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-microservice-k8s/services/customer/entity"
	"github.com/tusmasoma/go-microservice-k8s/services/customer/repository"
)

type CustomerUseCase interface {
	GetCustomer(ctx context.Context, id string) (*entity.Customer, error)
	ListCustomers(ctx context.Context) ([]entity.Customer, error)
	CreateCustomer(ctx context.Context, params *CreateCustomerParams) error
	UpdateCustomer(ctx context.Context, params *UpdateCustomerParams) error
	DeleteCustomer(ctx context.Context, id string) error
}

type customerUseCase struct {
	cr repository.CustomerRepository
}

func NewCustomerUsecase(cr repository.CustomerRepository) CustomerUseCase {
	return &customerUseCase{
		cr: cr,
	}
}

func (cuc *customerUseCase) GetCustomer(ctx context.Context, id string) (*entity.Customer, error) {
	customer, err := cuc.cr.Get(ctx, id)
	if err != nil {
		log.Error("failed to get customer", log.Ferror(err))
		return nil, err
	}
	return customer, nil
}

func (cuc *customerUseCase) ListCustomers(ctx context.Context) ([]entity.Customer, error) {
	customers, err := cuc.cr.List(ctx)
	if err != nil {
		log.Error("failed to list customers", log.Ferror(err))
		return nil, err
	}
	return customers, nil
}

type CreateCustomerParams struct {
	Name    string
	Email   string
	Street  string
	City    string
	Country string
}

func (cuc *customerUseCase) CreateCustomer(ctx context.Context, params *CreateCustomerParams) error {
	customer, err := entity.NewCustomer(
		params.Name,
		params.Email,
		params.Street,
		params.City,
		params.Country,
	)
	if err != nil {
		log.Error("failed to create customer", log.Ferror(err))
		return err
	}
	if err = cuc.cr.Create(ctx, *customer); err != nil {
		log.Error("failed to create customer", log.Ferror(err))
		return err
	}
	return nil
}

type UpdateCustomerParams struct {
	ID      string
	Name    string
	Email   string
	Street  string
	City    string
	Country string
}

func (cuc *customerUseCase) UpdateCustomer(ctx context.Context, params *UpdateCustomerParams) error {
	customer, err := cuc.cr.Get(ctx, params.ID)
	if err != nil {
		log.Error("failed to get customer", log.Ferror(err))
		return err
	}

	customer.Name = params.Name
	customer.Email = params.Email
	customer.Street = params.Street
	customer.City = params.City
	customer.Country = params.Country

	if err = cuc.cr.Update(ctx, *customer); err != nil {
		log.Error("failed to update customer", log.Ferror(err))
		return err
	}
	return nil
}

func (cuc *customerUseCase) DeleteCustomer(ctx context.Context, id string) error {
	if err := cuc.cr.Delete(ctx, id); err != nil {
		log.Error("failed to delete customer", log.Ferror(err))
		return err
	}
	return nil
}
