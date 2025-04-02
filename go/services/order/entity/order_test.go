package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEntity_NewOrder(t *testing.T) {
	t.Parallel()
	orderID := uuid.NewString()
	customerID := uuid.NewString()
	orderLines := OrderLines{
		&OrderLine{
			OrderLine: &order.OrderLine{
				Count:         1,
				CatalogItemId: uuid.NewString(),
			},
		},
	}
	now := time.Now()
	patterns := []struct {
		name string
		arg  struct {
			id         string
			customerID string
			orderDate  time.Time
			orderLines OrderLines
		}
		want struct {
			order *Order
			err   error
		}
	}{
		{
			name: "Success",
			arg: struct {
				id         string
				customerID string
				orderDate  time.Time
				orderLines OrderLines
			}{
				id:         orderID,
				customerID: customerID,
				orderDate:  now,
				orderLines: orderLines,
			},
			want: struct {
				order *Order
				err   error
			}{
				order: &Order{
					Order: order.Order{
						Id:         orderID,
						CustomerId: customerID,
						OrderDate:  timestamppb.New(now),
						OrderLines: orderLines.Proto(),
					},
				},
				err: nil,
			},
		},
		{
			name: "Fail: id is empty",
			arg: struct {
				id         string
				customerID string
				orderDate  time.Time
				orderLines OrderLines
			}{
				id:         "",
				customerID: customerID,
				orderDate:  now,
				orderLines: orderLines,
			},
			want: struct {
				order *Order
				err   error
			}{
				order: nil,
				err:   errors.New("id is required"),
			},
		},
		{
			name: "Fail: customerID is empty",
			arg: struct {
				id         string
				customerID string
				orderDate  time.Time
				orderLines OrderLines
			}{
				id:         orderID,
				customerID: "",
				orderDate:  now,
				orderLines: orderLines,
			},
			want: struct {
				order *Order
				err   error
			}{
				order: nil,
				err:   errors.New("customerID is required"),
			},
		},
		{
			name: "Fail: order lines is required",
			arg: struct {
				id         string
				customerID string
				orderDate  time.Time
				orderLines OrderLines
			}{
				id:         orderID,
				customerID: customerID,
				orderDate:  now,
				orderLines: nil,
			},
			want: struct {
				order *Order
				err   error
			}{
				order: nil,
				err:   errors.New("order lines is required"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			order, err := NewOrder(tt.arg.id, tt.arg.customerID, tt.arg.orderDate, tt.arg.orderLines)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.want.err)
			}
			if tt.want.err != nil {
				return
			}
			assert.Equal(t, tt.want.order.GetId(), order.GetId())
			assert.Equal(t, tt.want.order.GetCustomerId(), order.GetCustomerId())
			assert.Equal(t, tt.want.order.GetOrderDate(), order.GetOrderDate())
			assert.Equal(t, tt.want.order.GetOrderLines(), order.GetOrderLines())
		})
	}
}

func TestEntity_CreateOrder(t *testing.T) {
	t.Parallel()
	customerID := uuid.NewString()
	orderLines := OrderLines{
		&OrderLine{
			OrderLine: &order.OrderLine{
				Count:         1,
				CatalogItemId: uuid.NewString(),
			},
		},
	}
	patterns := []struct {
		name string
		arg  struct {
			customerID string
			orderLines OrderLines
		}
		want struct {
			order *Order
			err   error
		}
	}{
		{
			name: "Success",
			arg: struct {
				customerID string
				orderLines OrderLines
			}{
				customerID: customerID,
				orderLines: orderLines,
			},
			want: struct {
				order *Order
				err   error
			}{
				order: &Order{
					Order: order.Order{
						CustomerId: customerID,
						OrderLines: orderLines.Proto(),
					},
				},
				err: nil,
			},
		},
		{
			name: "Fail: customerID is empty",
			arg: struct {
				customerID string
				orderLines OrderLines
			}{
				customerID: "",
				orderLines: orderLines,
			},
			want: struct {
				order *Order
				err   error
			}{
				order: nil,
				err:   errors.New("customerID is required"),
			},
		},
		{
			name: "Fail: order lines is required",
			arg: struct {
				customerID string
				orderLines OrderLines
			}{
				customerID: customerID,
				orderLines: nil,
			},
			want: struct {
				order *Order
				err   error
			}{
				order: nil,
				err:   errors.New("order lines is required"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			order, err := CreateOrder(tt.arg.customerID, tt.arg.orderLines)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.want.err)
			}
			if tt.want.err != nil {
				return
			}
			assert.Equal(t, tt.want.order.GetCustomerId(), order.GetCustomerId())
			assert.Equal(t, tt.want.order.GetOrderLines(), order.GetOrderLines())
		})
	}
}

func Test_Order_Proto(t *testing.T) {
	now := time.Now()
	o := Order{
		Order: order.Order{
			Id:         "1",
			CustomerId: "1",
			OrderDate:  timestamppb.New(now),
			OrderLines: []*order.OrderLine{
				{
					Count:         1,
					CatalogItemId: "1",
				},
			},
		},
	}
	expect := &order.Order{
		Id:         "1",
		CustomerId: "1",
		OrderDate:  timestamppb.New(now),
		OrderLines: []*order.OrderLine{
			{
				Count:         1,
				CatalogItemId: "1",
			},
		},
	}
	assert.Equal(t, expect, o.Proto())
}

func Test_Orders_Proto(t *testing.T) {
	now := time.Now()
	o := Orders{
		&Order{
			order.Order{
				Id:         "1",
				CustomerId: "1",
				OrderDate:  timestamppb.New(now),
				OrderLines: []*order.OrderLine{
					{
						Count:         1,
						CatalogItemId: "1",
					},
				},
			},
		},
	}
	expect := []*order.Order{
		{
			Id:         "1",
			CustomerId: "1",
			OrderDate:  timestamppb.New(now),
			OrderLines: []*order.OrderLine{
				{
					Count:         1,
					CatalogItemId: "1",
				},
			},
		},
	}
	assert.Equal(t, expect, o.Proto())
}

func Test_OrderLine_Proto(t *testing.T) {
	ol := OrderLine{
		OrderLine: &order.OrderLine{
			Count:         1,
			CatalogItemId: "1",
		},
	}
	expect := &order.OrderLine{
		Count:         1,
		CatalogItemId: "1",
	}
	assert.Equal(t, expect, ol.Proto())
}

func Test_OrderLines_Proto(t *testing.T) {
	ols := OrderLines{
		&OrderLine{
			OrderLine: &order.OrderLine{
				Count:         1,
				CatalogItemId: "1",
			},
		},
	}
	expect := []*order.OrderLine{
		{
			Count:         1,
			CatalogItemId: "1",
		},
	}
	assert.Equal(t, expect, ols.Proto())
}
