package gateway

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tusmasoma/microservice-k8s-demo/customer/usecase"

	pb "github.com/tusmasoma/microservice-k8s-demo/customer/proto"
)

type CustomerHandler interface {
	ListCustomers(ctx context.Context, req *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error)
	CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error)
	UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error)
	DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error)
}

type customerHandler struct {
	cuc usecase.CustomerUseCase
	pb.UnimplementedCustomerServiceServer
}

func NewCustomerHandler(cuc usecase.CustomerUseCase) *customerHandler { //nolint:revive // This function is used in the test
	return &customerHandler{
		cuc: cuc,
	}
}

func (ch *customerHandler) ListCustomers(ctx context.Context, _ *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error) {
	customers, err := ch.cuc.ListCustomers(ctx)
	if err != nil {
		log.Error("Failed to list customers", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list customers")
	}

	var res []*pb.Customer
	for _, customer := range customers {
		res = append(res, &pb.Customer{
			Id:      customer.ID,
			Name:    customer.Name,
			Email:   customer.Email,
			Street:  customer.Street,
			City:    customer.City,
			Country: customer.Country,
		})
	}

	return &pb.ListCustomersResponse{
		Customers: res,
	}, nil
}

func (ch *customerHandler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	if !ch.isValidCreateCustomerRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	params := ch.convertCreateCustomerReqeuestToParams(req)
	if err := ch.cuc.CreateCustomer(ctx, params); err != nil {
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

func (ch *customerHandler) convertCreateCustomerReqeuestToParams(req *pb.CreateCustomerRequest) *usecase.CreateCustomerParams {
	return &usecase.CreateCustomerParams{
		Name:    req.GetName(),
		Email:   req.GetEmail(),
		Street:  req.GetStreet(),
		City:    req.GetCity(),
		Country: req.GetCountry(),
	}
}

func (ch *customerHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {
	if !ch.isValidUpdateCustomerRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	params := ch.convertUpdateCustomerReqeuestToParams(req)
	if err := ch.cuc.UpdateCustomer(ctx, params); err != nil {
		log.Error("Failed to update customer", log.Ferror(err))
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

func (ch *customerHandler) convertUpdateCustomerReqeuestToParams(req *pb.UpdateCustomerRequest) *usecase.UpdateCustomerParams {
	return &usecase.UpdateCustomerParams{
		ID:      req.GetId(),
		Name:    req.GetName(),
		Email:   req.GetEmail(),
		Street:  req.GetStreet(),
		City:    req.GetCity(),
		Country: req.GetCountry(),
	}
}

func (ch *customerHandler) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}

	if err := ch.cuc.DeleteCustomer(ctx, id); err != nil {
		log.Error("Failed to delete customer", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete customer")
	}

	return &pb.DeleteCustomerResponse{}, nil
}
