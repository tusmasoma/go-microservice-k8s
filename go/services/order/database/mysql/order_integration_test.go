package mysql

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-microservice-k8s/go/services/order/entity"
	pb "github.com/tusmasoma/go-microservice-k8s/proto/order"
)

func Test_OrderRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewOrder(db)
	customerID := uuid.NewString()
	lines := []*entity.OrderLine{
		{
			OrderLine: &pb.OrderLine{
				Count:         1,
				CatalogItemId: uuid.NewString(),
			},
		},
	}
	order, err := entity.CreateOrder(customerID, lines)
	if err != nil {
		return
	}
	// Create
	err = repo.Create(ctx, order)
	ValidateErr(t, err, nil)
	// Get
	gotOrder, err := repo.Get(ctx, order.GetId())
	ValidateErr(t, err, nil)
	if order.GetId() != gotOrder.GetId() {
		t.Errorf("unexpected order ID: want=%s, got=%s", order.GetId(), gotOrder.GetId())
	}
	// List
	gotOrders, err := repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gotOrders) != 1 {
		t.Errorf("got %d orders, want 1", len(gotOrders))
	}
	// Delete
	err = repo.Delete(ctx, order.GetId())
	ValidateErr(t, err, nil)
	_, err = repo.Get(ctx, order.GetId())
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
