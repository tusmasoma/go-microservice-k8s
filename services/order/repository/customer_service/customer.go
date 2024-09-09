package customerservice

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/services/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/services/order/repository"

	pb "github.com/tusmasoma/go-microservice-k8s/services/customer/proto"
)

type customerRepository struct {
	client pb.CustomerServiceClient
}

func NewCustomerRepository(client pb.CustomerServiceClient) repository.CustomerRepository {
	return &customerRepository{
		client,
	}
}

func (r *customerRepository) Get(ctx context.Context, id string) (*entity.Customer, error) {
	resp, err := r.client.GetCustomer(ctx, &pb.GetCustomerRequest{Id: id})
	if err != nil {
		return nil, err
	}

	customer := entity.Customer{
		ID:      resp.GetCustomer().GetId(),
		Name:    resp.GetCustomer().GetName(),
		Email:   resp.GetCustomer().GetEmail(),
		Street:  resp.GetCustomer().GetStreet(),
		City:    resp.GetCustomer().GetCity(),
		Country: resp.GetCustomer().GetCountry(),
	}

	return &customer, err
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
