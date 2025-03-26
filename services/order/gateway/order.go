package gateway

import (
	"context"

	pb "github.com/tusmasoma/go-microservice-k8s/proto/order"
	"github.com/tusmasoma/go-microservice-k8s/services/order/database"
	"github.com/tusmasoma/go-microservice-k8s/services/order/entity"
)

type OrderHandler interface {
	GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error)
	ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error)
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error)
}

type orderHandler struct {
	db *database.Database
	pb.UnimplementedOrderServiceServer
}

func NewOrderHandler(db *database.Database) pb.OrderServiceServer {
	return &orderHandler{
		db: db,
	}
}

func (oh *orderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := oh.db.Order.Get(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{
		Order: order.Proto(),
	}, nil
}

func (oh *orderHandler) ListOrders(ctx context.Context, _ *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := oh.db.Order.List(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListOrdersResponse{
		Orders: orders.Proto(),
	}, nil
}

func (oh *orderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	orderLines := make(entity.OrderLines, 0, len(req.GetOrderLines()))
	for _, ol := range req.GetOrderLines() {
		orderLines = append(orderLines, &entity.OrderLine{
			OrderLine: ol,
		})
	}
	order, err := entity.CreateOrder(req.GetCustomerId(), orderLines)
	if err != nil {
		return nil, err
	}
	if err := oh.db.Order.Create(ctx, order); err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{}, nil
}

func (oh *orderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if err := oh.db.Order.Delete(ctx, req.GetOrderId()); err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{}, nil
}
