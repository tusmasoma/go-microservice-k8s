package api

import (
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	"github.com/tusmasoma/go-microservice-k8s/services/order/database"
)

type orderService struct {
	db *database.Database
	order.UnimplementedOrderServiceServer
}

func NewOrderService(db *database.Database) order.OrderServiceServer {
	return &orderService{
		db: db,
	}
}
