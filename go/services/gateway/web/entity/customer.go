package entity

import (
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/proto/web"
)

type Customer struct {
	web.Customer
}

func (c *Customer) Proto() *web.Customer {
	if c == nil {
		return nil
	}
	return &c.Customer
}

func NewCustomer(customer *customer.Customer) *Customer {
	return &Customer{
		Customer: web.Customer{
			Id:      customer.GetId(),
			Name:    customer.GetName(),
			Email:   customer.GetEmail(),
			Street:  customer.GetStreet(),
			City:    customer.GetCity(),
			Country: customer.GetCountry(),
		},
	}
}

type Customers []*Customer

func NewCustomers(customers []*customer.Customer) Customers {
	ret := make(Customers, len(customers))
	for i := range customers {
		ret[i] = NewCustomer(customers[i])
	}
	return ret
}

func (cs Customers) Proto() []*web.Customer {
	customers := make([]*web.Customer, len(cs))
	for i, customer := range cs {
		customers[i] = customer.Proto()
	}
	return customers
}
