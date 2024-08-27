package gateway

import (
	"context"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/usecase"
)

type OrderHandler interface {
	GetOrderPageData(ctx context.Context, req *pb.GetOrderPageDataRequest) (*pb.GetOrderPageDataResponse, error)
}

type orderHandler struct {
	ouc usecase.OrderUseCase
	pb.UnimplementedOrderServiceServer
}

func NewOrderHandler(ouc usecase.OrderUseCase) pb.OrderServiceServer {
	return &orderHandler{
		ouc: ouc,
	}
}

func (oh *orderHandler) GetOrderPageData(ctx context.Context, _ *pb.GetOrderPageDataRequest) (*pb.GetOrderPageDataResponse, error) {
	customers, items, err := oh.ouc.GetOrderPageData(ctx)
	if err != nil {
		return nil, err
	}
	customerResponses := make([]*pb.Customer, 0, len(customers))
	for _, customer := range customers {
		customerResponses = append(customerResponses, &pb.Customer{
			Id:   customer.ID,
			Name: customer.Name,
		})
	}
	itemResponses := make([]*pb.CatalogItem, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, &pb.CatalogItem{
			Id:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		})
	}
	return &pb.GetOrderPageDataResponse{
		Customers: customerResponses,
		Items:     itemResponses,
	}, nil
}
