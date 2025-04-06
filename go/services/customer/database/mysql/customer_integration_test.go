package mysql

import (
	"context"
	"reflect"
	"testing"

	pt "github.com/tusmasoma/go-microservice-k8s/go/pkg/testing"
	"github.com/tusmasoma/go-microservice-k8s/go/services/customer/entity"
	pb "github.com/tusmasoma/go-microservice-k8s/proto/customer"
)

func Test_CustomerRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewCustomer(db)

	customer1, err := entity.CreateCustomer(
		"John Doe",
		"john.doe@example.com",
		"123 Maple Street",
		"Springfield",
		pb.CountryCode_JPN,
	)
	pt.ValidateErr(t, err, nil)
	customer2, err := entity.CreateCustomer(
		"Jane Smith",
		"jane.smith@example.com",
		"456 Oak Avenue",
		"Seattle",
		pb.CountryCode_JPN,
	)
	pt.ValidateErr(t, err, nil)

	// Create
	err = repo.Create(ctx, customer1)
	pt.ValidateErr(t, err, nil)
	err = repo.Create(ctx, customer2)
	pt.ValidateErr(t, err, nil)

	// Get
	gotCustomer, err := repo.Get(ctx, customer1.GetId())
	pt.ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gotCustomer, customer1) {
		t.Errorf("expected: %v, got: %v", customer1, gotCustomer)
	}

	// List
	gotCustomers, err := repo.List(ctx)
	pt.ValidateErr(t, err, nil)
	if len(gotCustomers) != 2 {
		t.Errorf("expected: 2, got: %d", len(gotCustomers))
	}

	// Update
	customer1.Name = "John Smith"
	err = repo.Update(ctx, customer1)
	pt.ValidateErr(t, err, nil)
	gotCustomer, err = repo.Get(ctx, customer1.GetId())
	pt.ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gotCustomer, customer1) {
		t.Errorf("expected: %v, got: %v", customer1, gotCustomer)
	}

	// Delete
	err = repo.Delete(ctx, customer1.GetId())
	pt.ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, customer1.GetId())
	if err == nil {
		t.Errorf("expected: error, got: nil")
	}
}
