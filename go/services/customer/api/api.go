package api

import (
	"github.com/tusmasoma/go-microservice-k8s/go/services/customer/database"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
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
