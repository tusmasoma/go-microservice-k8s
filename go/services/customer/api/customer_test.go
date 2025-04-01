package api

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tusmasoma/go-microservice-k8s/go/services/customer/database"
	dbmock "github.com/tusmasoma/go-microservice-k8s/go/services/customer/database/mock"
	"github.com/tusmasoma/go-microservice-k8s/go/services/customer/entity"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
)

func Test_GetCustomer(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCustomer(ctrl)
	gomock.InOrder(
		db.EXPECT().Get(ctx, id).Return(nil, errmock),
		db.EXPECT().Get(ctx, id).Return(
			&entity.Customer{
				Customer: customer.Customer{
					Id:   id,
					Name: "person",
				},
			}, nil),
	)
	service := &customerService{
		db: &database.Database{
			Customer: db,
		},
	}
	// bad request
	req := &customer.GetCustomerRequest{}
	_, err := service.GetCustomer(ctx, req)
	assert.Error(t, err)
	// db error
	req.Id = id
	resp, err := service.GetCustomer(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.GetCustomer(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_ListCustomers(t *testing.T) {
	t.Parallel()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCustomer(ctrl)
	gomock.InOrder(
		db.EXPECT().List(ctx).Return(nil, errmock),
		db.EXPECT().List(ctx).Return(
			entity.Customers{
				&entity.Customer{
					Customer: customer.Customer{
						Id:   uuid.NewString(),
						Name: "person",
					},
				},
			}, nil),
	)
	service := &customerService{
		db: &database.Database{
			Customer: db,
		},
	}
	// db error
	req := &customer.ListCustomersRequest{}
	resp, err := service.ListCustomers(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.ListCustomers(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_CreateCustomer(t *testing.T) {
	t.Parallel()
	user := &entity.Customer{
		Customer: customer.Customer{
			Id:      uuid.NewString(),
			Name:    "John Doe",
			Email:   "john.doe@example.com",
			Street:  "1600 Pennsylvania Avenue NW",
			City:    "Washington",
			Country: "USA",
		},
	}
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCustomer(ctrl)
	gomock.InOrder(
		db.EXPECT().Create(
			gomock.Any(),
			gomock.Any(),
		).Return(errmock),
		db.EXPECT().Create(
			gomock.Any(),
			gomock.Any(),
		).Do(func(_ context.Context, user *entity.Customer) {
			if user.GetName() != "John Doe" {
				t.Errorf("unexpected Name: got %v, want %v", user.GetName(), "John Doe")
			}
			if user.GetEmail() != "john.doe@example.com" {
				t.Errorf("unexpected Email: got %v, want %v", user.GetEmail(), "john.doe@example.com")
			}
			if user.GetStreet() != "1600 Pennsylvania Avenue NW" {
				t.Errorf("unexpected Street: got %v, want %v", user.GetStreet(), "1600 Pennsylvania Avenue NW")
			}
			if user.GetCity() != "Washington" {
				t.Errorf("unexpected City: got %v, want %v", user.GetCity(), "Washington")
			}
			if user.GetCountry() != "USA" {
				t.Errorf("unexpected Country: got %v, want %v", user.GetCountry(), "USA")
			}
		}).Return(nil),
	)
	service := &customerService{
		db: &database.Database{
			Customer: db,
		},
	}
	// bad request
	req := &customer.CreateCustomerRequest{}
	_, err := service.CreateCustomer(ctx, req)
	assert.Error(t, err)
	// db error
	req.Name = user.GetName()
	req.Email = user.GetEmail()
	req.Street = user.GetStreet()
	req.City = user.GetCity()
	req.Country = user.GetCountry()
	resp, err := service.CreateCustomer(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.CreateCustomer(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_UpdateCustomer(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	user := &entity.Customer{
		Customer: customer.Customer{
			Id:      id,
			Name:    "John Doe",
			Email:   "john.doe@example.com",
			Street:  "1600 Pennsylvania Avenue NW",
			City:    "Washington",
			Country: "USA",
		},
	}
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCustomer(ctrl)
	gomock.InOrder(
		db.EXPECT().Get(ctx, id).Return(nil, errmock),
		db.EXPECT().Get(ctx, id).Return(user, nil),
		db.EXPECT().Update(
			gomock.Any(),
			gomock.Any(),
		).Return(errmock),
		db.EXPECT().Get(ctx, id).Return(user, nil),
		db.EXPECT().Update(
			gomock.Any(),
			gomock.Any(),
		).Do(func(_ context.Context, user *entity.Customer) {
			if user.GetName() != "John Doe" {
				t.Errorf("unexpected Name: got %v, want %v", user.GetName(), "John Doe")
			}
			if user.GetEmail() != "john.doe@example.com" {
				t.Errorf("unexpected Email: got %v, want %v", user.GetEmail(), "john.doe@example.com")
			}
			if user.GetStreet() != "1600 Pennsylvania Avenue NW" {
				t.Errorf("unexpected Street: got %v, want %v", user.GetStreet(), "1600 Pennsylvania Avenue NW")
			}
			if user.GetCity() != "Washington" {
				t.Errorf("unexpected City: got %v, want %v", user.GetCity(), "Washington")
			}
			if user.GetCountry() != "USA" {
				t.Errorf("unexpected Country: got %v, want %v", user.GetCountry(), "USA")
			}
		}).Return(nil),
	)
	service := &customerService{
		db: &database.Database{
			Customer: db,
		},
	}
	// bad request
	req := &customer.UpdateCustomerRequest{}
	_, err := service.UpdateCustomer(ctx, req)
	assert.Error(t, err)
	// db error (Get)
	req.Id = user.GetId()
	req.Name = user.GetName()
	req.Email = user.GetEmail()
	req.Street = user.GetStreet()
	req.City = user.GetCity()
	req.Country = user.GetCountry()
	resp, err := service.UpdateCustomer(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// db error (Update)
	resp, err = service.UpdateCustomer(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.UpdateCustomer(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_DeleteCustomer(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCustomer(ctrl)
	gomock.InOrder(
		db.EXPECT().Delete(ctx, id).Return(errmock),
		db.EXPECT().Delete(ctx, id).Return(nil),
	)
	service := &customerService{
		db: &database.Database{
			Customer: db,
		},
	}
	// bad request
	req := &customer.DeleteCustomerRequest{}
	_, err := service.DeleteCustomer(ctx, req)
	assert.Error(t, err)
	// db error
	req.Id = id
	resp, err := service.DeleteCustomer(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.DeleteCustomer(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
