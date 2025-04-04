package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
)

type Customer struct {
	customer.Customer
}

func (c *Customer) Proto() *customer.Customer {
	if c == nil {
		return nil
	}
	return &c.Customer
}

func NewCustomer(id, name, email, street, city string, countryCode customer.CountryCode) (*Customer, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if street == "" {
		return nil, errors.New("street is required")
	}
	if city == "" {
		return nil, errors.New("city is required")
	}
	if countryCode == customer.CountryCode_COUNTRY_CODE_UNSPECIFIED {
		return nil, errors.New("countryCode is required")
	}
	return &Customer{
		Customer: customer.Customer{
			Id:          id,
			Name:        name,
			Email:       email,
			Street:      street,
			City:        city,
			CountryCode: countryCode,
		},
	}, nil
}

func CreateCustomer(name, email, street, city string, countryCode customer.CountryCode) (*Customer, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if street == "" {
		return nil, errors.New("street is required")
	}
	if city == "" {
		return nil, errors.New("city is required")
	}
	if countryCode == customer.CountryCode_COUNTRY_CODE_UNSPECIFIED {
		return nil, errors.New("countryCode is required")
	}
	return &Customer{
		Customer: customer.Customer{
			Id:          uuid.NewString(),
			Name:        name,
			Email:       email,
			Street:      street,
			City:        city,
			CountryCode: countryCode,
		},
	}, nil
}

type Customers []*Customer

func (cs Customers) Proto() []*customer.Customer {
	customers := make([]*customer.Customer, len(cs))
	for i, customer := range cs {
		customers[i] = customer.Proto()
	}
	return customers
}
