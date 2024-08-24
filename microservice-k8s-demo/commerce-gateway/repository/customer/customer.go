package customer

import (
	"context"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/proto"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/commerce-gateway/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/commerce-gateway/repository"
)

type customerRepository struct {
	client pb.CustomerServiceClient
}

func NewCustomerRepository(client pb.CustomerServiceClient) repository.CustomerRepository {
	return &customerRepository{
		client,
	}
}

func (r *customerRepository) List(ctx context.Context) ([]entity.Customer, error) {
	resp, err := r.client.ListCustomers(ctx, &pb.ListCustomersRequest{})
	if err != nil {
		return nil, err
	}

	var customers []entity.Customer
	for _, c := range resp.GetCustomers() {
		customers = append(customers, entity.Customer{
			ID:      c.GetId(),
			Name:    c.GetName(),
			Email:   c.GetEmail(),
			Street:  c.GetStreet(),
			City:    c.GetCity(),
			Country: c.GetCountry(),
		})
	}

	return customers, nil
}

func (r *customerRepository) Create(ctx context.Context, customer entity.Customer) error {
	if _, err := r.client.CreateCustomer(ctx, &pb.CreateCustomerRequest{
		Name:    customer.Name,
		Email:   customer.Email,
		Street:  customer.Street,
		City:    customer.City,
		Country: customer.Country,
	}); err != nil {
		return err
	}
	return nil
}

func (r *customerRepository) Update(ctx context.Context, customer entity.Customer) error {
	if _, err := r.client.UpdateCustomer(ctx, &pb.UpdateCustomerRequest{
		Id:      customer.ID,
		Name:    customer.Name,
		Email:   customer.Email,
		Street:  customer.Street,
		City:    customer.City,
		Country: customer.Country,
	}); err != nil {
		return err
	}
	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.client.DeleteCustomer(ctx, &pb.DeleteCustomerRequest{Id: id}); err != nil {
		return err
	}
	return nil
}
