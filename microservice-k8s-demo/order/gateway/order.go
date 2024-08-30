package gateway

import (
	"context"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/usecase"
)

type OrderHandler interface {
	GetOrderCreationResources(ctx context.Context, req *pb.GetOrderCreationResourcesRequest) (*pb.GetOrderCreationResourcesResponse, error)
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
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

func (oh *orderHandler) GetOrderCreationResources(ctx context.Context, _ *pb.GetOrderCreationResourcesRequest) (*pb.GetOrderCreationResourcesResponse, error) {
	customers, items, err := oh.ouc.GetOrderCreationResources(ctx)
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
	return &pb.GetOrderCreationResourcesResponse{
		Customers: customerResponses,
		Items:     itemResponses,
	}, nil
}

func (oh *orderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	orderLines := make([]struct {
		CatalogItemID string
		Count         int
	}, 0, len(req.GetOrderLines()))

	for _, ol := range req.GetOrderLines() {
		orderLines = append(orderLines, struct {
			CatalogItemID string
			Count         int
		}{
			CatalogItemID: ol.GetItemId(),
			Count:         int(ol.GetCount()),
		})
	}
	if err := oh.ouc.CreateOrder(ctx, &usecase.CreateOrderParams{
		CustomerID: req.GetCustomerId(),
		OrderLine:  orderLines,
	}); err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{}, nil
}
