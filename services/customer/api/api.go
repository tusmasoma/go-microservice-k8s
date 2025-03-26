package api

import (
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/services/customer/database"
)

type customerService struct {
	db *database.Database
	customer.UnimplementedCustomerServiceServer
}

func NewCustomerService(db *database.Database) customer.CustomerServiceServer {
	return &customerService{
		db: db,
	}
}
