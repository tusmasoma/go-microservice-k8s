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

	orderLine1 := entity.OrderLine{
		CatalogItemID: uuid.New().String(),
		Count:         1,
	}
	orderLine2 := entity.OrderLine{
		CatalogItemID: uuid.New().String(),
		Count:         5,
	}
	orderLine3 := entity.OrderLine{
		CatalogItemID: uuid.New().String(),
		Count:         2,
	}

	order1 := entity.Order{
		ID:         uuid.New().String(),
		CustomerID: uuid.New().String(),
		OrderDate:  time.Now(),
		OrderLines: []entity.OrderLine{orderLine1, orderLine2},
	}

	order2 := entity.Order{
		ID:         uuid.New().String(),
		CustomerID: uuid.New().String(),
		OrderDate:  time.Now(),
		OrderLines: []entity.OrderLine{orderLine2, orderLine3},
	}

	// Create
	err := repo.Create(ctx, order1)
	ValidateErr(t, err, nil)
	err = repo.Create(ctx, order2)
	ValidateErr(t, err, nil)

	// Get
	gotOrder1, err := repo.Get(ctx, order1.ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(order1, *gotOrder1,
		cmpopts.IgnoreFields(entity.Order{}, "OrderDate"),
		cmp.Comparer(func(x, y entity.Order) bool {
			return cmp.Equal(x.OrderLines, y.OrderLines,
				cmpopts.SortSlices(func(a, b entity.OrderLine) bool {
					return a.CatalogItemID < b.CatalogItemID
				}))
		})); len(d) != 0 {
		t.Errorf("differs: (-want +got)\n%s", d)
	}

	// List
	gotOrders, err := repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gotOrders) != 2 {
		t.Errorf("want: 2, got: %d", len(gotOrders))
	}

	// Delete
	err = repo.Delete(ctx, order1.ID)
	ValidateErr(t, err, nil)

	gotOrders, err = repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gotOrders) != 1 {
		t.Errorf("want: 1, got: %d", len(gotOrders))
	}
}
