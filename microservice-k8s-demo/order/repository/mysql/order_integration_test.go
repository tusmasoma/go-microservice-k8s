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
	orderRepo := NewOrderRepository(db)
	orderLienRepo := NewOrderLineRepository(db)

	order1ID := uuid.New().String()
	order2ID := uuid.New().String()

	t.Skip()

	orderLine1 := entity.OrderLineModel{
		OrderID:       order1ID,
		CatalogItemID: uuid.New().String(),
		Count:         1,
	}
	orderLine2 := entity.OrderLineModel{
		OrderID:       order1ID,
		CatalogItemID: uuid.New().String(),
		Count:         5,
	}
	orderLine3 := entity.OrderLineModel{
		OrderID:       order2ID,
		CatalogItemID: uuid.New().String(),
		Count:         2,
	}

	order1 := entity.OrderModel{
		ID:         order1ID,
		CustomerID: uuid.New().String(),
		OrderDate:  time.Now(),
	}

	order2 := entity.OrderModel{
		ID:         order2ID,
		CustomerID: uuid.New().String(),
		OrderDate:  time.Now(),
	}

	// Create
	err := orderRepo.Create(ctx, order1)
	ValidateErr(t, err, nil)
	err = orderRepo.Create(ctx, order2)
	ValidateErr(t, err, nil)

	err = orderLienRepo.Create(ctx, orderLine3)
	ValidateErr(t, err, nil)

	err = orderLienRepo.BatchCreate(ctx, []entity.OrderLineModel{orderLine1, orderLine2})
	ValidateErr(t, err, nil)

	// Get
	gotOrder1, err := orderRepo.Get(ctx, order1.ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(order1, *gotOrder1, cmpopts.IgnoreFields(entity.OrderModel{}, "OrderDate")); len(d) != 0 {
		t.Errorf("differs: (-want +got)\n%s", d)
	}

	// List
	gotOrders, err := orderRepo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gotOrders) != 2 {
		t.Errorf("want: 2, got: %d", len(gotOrders))
	}

	gotOrderLines, err := orderLienRepo.List(ctx, order1.ID)
	ValidateErr(t, err, nil)
	if len(gotOrderLines) != 2 {
		t.Errorf("want: 2, got: %d", len(gotOrderLines))
	}

	// Delete
	err = orderLienRepo.Delete(ctx, orderLine3.OrderID, orderLine3.CatalogItemID)
	ValidateErr(t, err, nil)

	err = orderLienRepo.BatchDelete(ctx, order1.ID)
	ValidateErr(t, err, nil)

	err = orderRepo.Delete(ctx, order1.ID)
	ValidateErr(t, err, nil)

	gotOrders, err = orderRepo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gotOrders) != 1 {
		t.Errorf("want: 1, got: %d", len(gotOrders))
	}
}
