package gateway

import (
	"context"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler interface {
	ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error)
	GetOrderCreationResources(ctx context.Context, req *pb.GetOrderCreationResourcesRequest) (*pb.GetOrderCreationResourcesResponse, error)
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error)
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

func (oh *orderHandler) ListOrders(ctx context.Context, _ *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orderDetails, err := oh.ouc.ListOrders(ctx)
	if err != nil {
		return nil, err
	}
	orderResponses := make([]*pb.Order, 0, len(orderDetails))
	for _, od := range orderDetails {
		orderLines := make([]*pb.OrderLine, 0, len(od.OrderLines))
		for _, ol := range od.OrderLines {
			orderLines = append(orderLines, &pb.OrderLine{
				Item: &pb.CatalogItem{
					Id:    ol.CatalogItem.ID,
					Name:  ol.CatalogItem.Name,
					Price: ol.CatalogItem.Price,
				},
				Count: int32(ol.Count),
			})
		}

		orderResponses = append(orderResponses, &pb.Order{
			Id: od.Order.ID,
			Customer: &pb.Customer{
				Id:      od.Customer.ID,
				Name:    od.Customer.Name,
				Email:   od.Customer.Email,
				Street:  od.Customer.Street,
				City:    od.Customer.City,
				Country: od.Customer.Country,
			},
			OrderDate:  timestamppb.New(od.Order.OrderDate),
			OrderLines: orderLines,
			TotalPrice: od.Order.TotalPrice,
		})
	}
	return &pb.ListOrdersResponse{
		Orders: orderResponses,
	}, nil
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
			CatalogItemID: ol.GetItem().GetId(),
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

func (oh *orderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if err := oh.ouc.DeleteOrder(ctx, req.GetOrderId()); err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{}, nil
}
