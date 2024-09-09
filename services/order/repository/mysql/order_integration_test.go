package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/tusmasoma/go-microservice-k8s/services/order/entity"
)

func Test_OrderRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewOrderRepository(db)

	order := entity.Order{
		ID:         uuid.New().String(),
		CustomerID: uuid.New().String(),
		OrderDate:  time.Now(),
		OrderLines: []*entity.OrderLine{
			{
				CatalogItemID: uuid.New().String(),
				Count:         1,
			},
		},
	}

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
