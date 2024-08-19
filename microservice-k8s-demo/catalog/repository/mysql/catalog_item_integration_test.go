package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/entity"
)

func Test_CatalogItemRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewCatalogItemRepository(db)

	item1, err := entity.NewCatalogItem(
		"item1",
		100,
	)
	ValidateErr(t, err, nil)
	item2, err := entity.NewCatalogItem(
		"item2",
		200,
	)
	ValidateErr(t, err, nil)

	// Create
	err = repo.Create(ctx, *item1)
	ValidateErr(t, err, nil)
	err = repo.Create(ctx, *item2)
	ValidateErr(t, err, nil)

	// Get
	gotItem, err := repo.Get(ctx, item1.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gotItem, item1) {
		t.Errorf("want: %v, got: %v", item1, gotItem)
	}

	// List
	gotItems, err := repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gotItems) != 2 {
		t.Errorf("want: 2, got: %d", len(gotItems))
	}

	// Update
	item1.Name = "item1-updated"
	item1.Price = 150
	err = repo.Update(ctx, *item1)
	ValidateErr(t, err, nil)

	gotItem, err = repo.Get(ctx, item1.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gotItem, item1) {
		t.Errorf("want: %v, got: %v", item1, gotItem)
	}

	// Delete
	err = repo.Delete(ctx, item1.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, item1.ID)
	if err == nil {
		t.Errorf("want: error, got: nil")
	}
}
