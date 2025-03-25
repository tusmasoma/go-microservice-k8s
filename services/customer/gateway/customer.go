package gateway

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tusmasoma/go-microservice-k8s/services/customer/database"
	"github.com/tusmasoma/go-microservice-k8s/services/customer/entity"

	pb "github.com/tusmasoma/go-microservice-k8s/proto/customer"
)

type CustomerHandler interface {
	GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error)
	ListCustomers(ctx context.Context, req *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error)
	CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error)
	UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error)
	DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error)
}

type customerHandler struct {
	db *database.Database
	pb.UnimplementedCustomerServiceServer
}

func NewCustomerHandler(db *database.Database) pb.CustomerServiceServer {
	return &customerHandler{
		db: db,
	}
}

func (ch *customerHandler) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	customer, err := ch.db.Customer.Get(ctx, id)
	if err != nil {
		log.Error("Failed to get customer", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to get customer")
	}
	return &pb.GetCustomerResponse{
		Customer: customer.Proto(),
	}, nil
}

func (ch *customerHandler) ListCustomers(ctx context.Context, _ *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error) {
	customers, err := ch.db.Customer.List(ctx)
	if err != nil {
		log.Error("Failed to list customers", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list customers")
	}
	return &pb.ListCustomersResponse{
		Customers: customers.Proto(),
	}, nil
}

func (ch *customerHandler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	if !ch.isValidCreateCustomerRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	customer, err := entity.CreateCustomer(
		req.GetName(),
		req.GetEmail(),
		req.GetStreet(),
		req.GetCity(),
		req.GetCountry(),
	)
	if err != nil {
		return nil, err
	}
	if err := ch.db.Customer.Create(ctx, customer); err != nil {
		log.Error("Failed to create customer", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to create customer")
	}
	return &pb.CreateCustomerResponse{}, nil
}

func (ch *customerHandler) isValidCreateCustomerRequest(req *pb.CreateCustomerRequest) bool {
	if req.GetName() == "" ||
		req.GetEmail() == "" ||
		req.GetStreet() == "" ||
		req.GetCity() == "" ||
		req.GetCountry() == "" {
		log.Warn(
			"Invalid request",
			log.Fstring("name", req.GetName()),
			log.Fstring("email", req.GetEmail()),
			log.Fstring("street", req.GetStreet()),
			log.Fstring("city", req.GetCity()),
			log.Fstring("country", req.GetCountry()),
		)
		return false
	}
	return true
}

func (ch *customerHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {
	if !ch.isValidUpdateCustomerRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	customer, err := ch.db.Customer.Get(ctx, req.GetId())
	if err != nil {
		log.Error("failed to get customer", log.Ferror(err))
		return nil, err
	}
	customer.Name = req.GetName()
	customer.Email = req.GetEmail()
	customer.Street = req.GetStreet()
	customer.City = req.GetCity()
	customer.Country = req.GetCountry()
	if err = ch.db.Customer.Update(ctx, customer); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update customer")
	}
	return &pb.UpdateCustomerResponse{}, nil
}

func (ch *customerHandler) isValidUpdateCustomerRequest(req *pb.UpdateCustomerRequest) bool {
	if req.GetId() == "" ||
		req.GetName() == "" ||
		req.GetEmail() == "" ||
		req.GetStreet() == "" ||
		req.GetCity() == "" ||
		req.GetCountry() == "" {
		log.Warn("Invalid request body: %v", req)
		return false
	}
	return true
}

func (ch *customerHandler) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	if err := ch.db.Customer.Delete(ctx, id); err != nil {
		log.Error("Failed to delete customer", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete customer")
	}
	return &pb.DeleteCustomerResponse{}, nil
}
