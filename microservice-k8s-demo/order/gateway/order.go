package gateway

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/usecase"
)

type OrderHandler interface {
	GetOrderPageData(ctx context.Context, req *pb.GetOrderPageDataRequest) (*pb.GetOrderPageDataResponse, error)
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

func (oh *orderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	orderLines := make([]entity.OrderLine, 0, len(req.GetOrderLines()))
	for _, ol := range req.GetOrderLines() {
		orderLines = append(orderLines, entity.OrderLine{
			Count:         int(ol.GetCount()),
			CatalogItemID: ol.GetItemId(),
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
