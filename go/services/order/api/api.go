package api

import (
	"github.com/tusmasoma/go-microservice-k8s/go/services/order/database"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
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
