package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tusmasoma/go-microservice-k8s/go/services/customer/entity"

	pb "github.com/tusmasoma/go-microservice-k8s/proto/customer"
)

func (c *customerService) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	customer, err := c.db.Customer.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get customer")
	}
	return &pb.GetCustomerResponse{
		Customer: customer.Proto(),
	}, nil
}

func (c *customerService) ListCustomers(ctx context.Context, _ *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error) {
	customers, err := c.db.Customer.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list customers")
	}
	return &pb.ListCustomersResponse{
		Customers: customers.Proto(),
	}, nil
}

func (c *customerService) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	if !c.isValidCreateCustomerRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	customer, err := entity.CreateCustomer(
		req.GetName(),
		req.GetEmail(),
		req.GetStreet(),
		req.GetCity(),
		req.GetCountryCode(),
	)
	if err != nil {
		return nil, err
	}
	if err := c.db.Customer.Create(ctx, customer); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create customer")
	}
	return &pb.CreateCustomerResponse{}, nil
}

func (c *customerService) isValidCreateCustomerRequest(req *pb.CreateCustomerRequest) bool {
	if req.GetName() == "" ||
		req.GetEmail() == "" ||
		req.GetStreet() == "" ||
		req.GetCity() == "" ||
		req.GetCountryCode() == pb.CountryCode_COUNTRY_CODE_UNSPECIFIED {
		return false
	}
	return true
}

func (c *customerService) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {
	if !c.isValidUpdateCustomerRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	customer, err := c.db.Customer.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	customer.Name = req.GetName()
	customer.Email = req.GetEmail()
	customer.Street = req.GetStreet()
	customer.City = req.GetCity()
	customer.CountryCode = req.GetCountryCode()
	if err = c.db.Customer.Update(ctx, customer); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update customer")
	}
	return &pb.UpdateCustomerResponse{}, nil
}

func (c *customerService) isValidUpdateCustomerRequest(req *pb.UpdateCustomerRequest) bool {
	if req.GetId() == "" ||
		req.GetName() == "" ||
		req.GetEmail() == "" ||
		req.GetStreet() == "" ||
		req.GetCity() == "" ||
		req.GetCountryCode() == pb.CountryCode_COUNTRY_CODE_UNSPECIFIED {
		return false
	}
	return true
}

func (c *customerService) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	if err := c.db.Customer.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete customer")
	}
	return &pb.DeleteCustomerResponse{}, nil
}
