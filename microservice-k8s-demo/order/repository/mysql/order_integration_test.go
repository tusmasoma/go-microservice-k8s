package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
)

func Test_OrderRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewOrderRepository(db)

	order := entity.Order{
		ID: uuid.New().String(),
		Customer: entity.Customer{
			ID: uuid.New().String(),
			// Name:    "John Doe",
			// Email:   "john.doe@example.com",
			// Street:  "1600 Pennsylvania Avenue NW",
			// City:    "Washington",
			// Country: "USA",
		},
		OrderDate: time.Now(),
		OrderLines: []entity.OrderLine{
			{
				CatalogItem: entity.CatalogItem{
					ID: uuid.New().String(),
					// Name:  "item1",
					// Price: 100,
				},
				Count: 1,
			},
		},
	}
	order.TotalPrice = order.GetTotalPrice()

	// Create
	err := repo.Create(ctx, order)
	ValidateErr(t, err, nil)

	// Get
	gotOrder, err := repo.Get(ctx, order.ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(order, *gotOrder, cmpopts.IgnoreFields(entity.Order{}, "OrderDate")); len(d) != 0 {
		t.Errorf("differs: (-want +got)\n%s", d)
	}

	// List
	gotOrders, err := repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gotOrders) != 1 {
		t.Errorf("got %d orders, want 1", len(gotOrders))
	}

	// Delete
	err = repo.Delete(ctx, order.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, order.ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
