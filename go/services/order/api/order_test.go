package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tusmasoma/go-microservice-k8s/go/services/order/database"
	dbmock "github.com/tusmasoma/go-microservice-k8s/go/services/order/database/mock"
	"github.com/tusmasoma/go-microservice-k8s/go/services/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_GetOrder(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockOrder(ctrl)
	gomock.InOrder(
		db.EXPECT().Get(ctx, id).Return(nil, errmock),
		db.EXPECT().Get(ctx, id).Return(
			&entity.Order{
				Order: order.Order{
					Id:         "1",
					CustomerId: "1",
					OrderDate:  timestamppb.New(time.Now()),
					OrderLines: []*order.OrderLine{
						{
							Quantity:      1,
							CatalogItemId: "1",
						},
					},
				},
			}, nil),
	)
	service := &orderService{
		db: &database.Database{
			Order: db,
		},
	}
	// bad request
	req := &order.GetOrderRequest{}
	_, err := service.GetOrder(ctx, req)
	assert.Error(t, err)
	// db error
	req.Id = id
	resp, err := service.GetOrder(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.GetOrder(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_ListOrders(t *testing.T) {
	t.Parallel()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockOrder(ctrl)
	gomock.InOrder(
		db.EXPECT().List(ctx).Return(nil, errmock),
		db.EXPECT().List(ctx).Return(
			entity.Orders{
				&entity.Order{
					Order: order.Order{
						Id:         "1",
						CustomerId: "1",
						OrderDate:  timestamppb.New(time.Now()),
						OrderLines: []*order.OrderLine{
							{
								Quantity:      1,
								CatalogItemId: "1",
							},
						},
					},
				},
			}, nil),
	)
	service := &orderService{
		db: &database.Database{
			Order: db,
		},
	}
	// db error
	req := &order.ListOrdersRequest{}
	resp, err := service.ListOrders(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.ListOrders(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_CreateOrder(t *testing.T) {
	t.Parallel()
	customerID := uuid.NewString()
	lines := entity.OrderLines{
		&entity.OrderLine{
			OrderLine: &order.OrderLine{
				Quantity:      1,
				CatalogItemId: "1",
			},
		},
	}
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockOrder(ctrl)
	gomock.InOrder(
		db.EXPECT().Create(ctx, gomock.Any()).Return(errmock),
		db.EXPECT().Create(ctx, gomock.Any()).Do(func(_ context.Context, order *entity.Order) {
			if order.GetCustomerId() != customerID {
				t.Errorf("unexpected CustomerID: got %v, want %v", order.GetCustomerId(), customerID)
			}
			if len(order.GetOrderLines()) != len(lines) {
				t.Errorf("unexpected OrderLines: got %v, want %v", order.GetOrderLines(), lines)
			}
		}).Return(nil),
	)
	service := &orderService{
		db: &database.Database{
			Order: db,
		},
	}
	// db error
	req := &order.CreateOrderRequest{}
	req.CustomerId = customerID
	req.OrderLines = lines.Proto()
	resp, err := service.CreateOrder(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.CreateOrder(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_DeleteOrder(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockOrder(ctrl)
	gomock.InOrder(
		db.EXPECT().Delete(ctx, id).Return(errmock),
		db.EXPECT().Delete(ctx, id).Return(nil),
	)
	service := &orderService{
		db: &database.Database{
			Order: db,
		},
	}
	// bad request
	req := &order.DeleteOrderRequest{}
	_, err := service.DeleteOrder(ctx, req)
	assert.Error(t, err)
	// db error
	req.Id = id
	resp, err := service.DeleteOrder(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.DeleteOrder(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
