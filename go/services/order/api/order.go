package api

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/go/services/order/entity"
	pb "github.com/tusmasoma/go-microservice-k8s/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *orderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	order, err := o.db.Order.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{
		Order: order.Proto(),
	}, nil
}

func (o *orderService) ListOrders(ctx context.Context, _ *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := o.db.Order.List(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListOrdersResponse{
		Orders: orders.Proto(),
	}, nil
}

func (o *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
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
	if err := o.db.Order.Create(ctx, order); err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{}, nil
}

func (o *orderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	if err := o.db.Order.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{}, nil
}
