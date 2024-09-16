package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Customer struct {
	ID      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Email   string `json:"email" db:"email"`
	Street  string `json:"street" db:"street"`
	City    string `json:"city" db:"city"`
	Country string `json:"country" db:"country"`
}

func NewCustomer(id, name, email, street, city, country string) (*Customer, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if name == "" {
		log.Error("name is required")
		return nil, errors.New("name is required")
	}
	if email == "" {
		log.Error("email is required")
		return nil, errors.New("email is required")
	}
	if street == "" {
		log.Error("street is required")
		return nil, errors.New("street is required")
	}
	if city == "" {
		log.Error("city is required")
		return nil, errors.New("city is required")
	}
	if country == "" {
		log.Error("country is required")
		return nil, errors.New("country is required")
	}
	return &Customer{
		ID:      id,
		Name:    name,
		Email:   email,
		Street:  street,
		City:    city,
		Country: country,
	}, nil
}
