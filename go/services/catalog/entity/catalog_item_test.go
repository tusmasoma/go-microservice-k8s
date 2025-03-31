package entity

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
)

func TestEntity_NewCatalogItem(t *testing.T) {
	t.Parallel()
	catalogID := uuid.New().String()
	patterns := []struct {
		name string
		arg  struct {
			id    string
			name  string
			price float64
		}
		want struct {
			item *CatalogItem
			err  error
		}
	}{
		{
			name: "Success",
			arg: struct {
				id    string
				name  string
				price float64
			}{
				id:    catalogID,
				name:  "item",
				price: 100,
			},
			want: struct {
				item *CatalogItem
				err  error
			}{
				item: &CatalogItem{
					CatalogItem: catalog.CatalogItem{
						Id:    catalogID,
						Name:  "item",
						Price: 100,
					},
				},
				err: nil,
			},
		},
		{
			name: "success: id is empty",
			arg: struct {
				id    string
				name  string
				price float64
			}{
				name:  "item",
				price: 100,
			},
			want: struct {
				item *CatalogItem
				err  error
			}{
				item: nil,
				err:  errors.New("id is required"),
			},
		},
		{
			name: "Fail: name is empty",
			arg: struct {
				id    string
				name  string
				price float64
			}{
				id:    catalogID,
				name:  "",
				price: 100,
			},
			want: struct {
				item *CatalogItem
				err  error
			}{
				item: nil,
				err:  errors.New("name is required"),
			},
		},
		{
			name: "Fail: price is less than 0",
			arg: struct {
				id    string
				name  string
				price float64
			}{
				id:    catalogID,
				name:  "item",
				price: -1,
			},
			want: struct {
				item *CatalogItem
				err  error
			}{
				item: nil,
				err:  errors.New("price must be greater than 0"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			item, err := NewCatalogItem(tt.arg.id, tt.arg.name, tt.arg.price)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewCatalogItem() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewCatalogItem() error = %v, wantErr %v", err, tt.want.err)
			}
			if tt.want.err != nil {
				return
			}
			assert.Equal(t, item.GetId(), tt.want.item.GetId())
			assert.Equal(t, item.GetName(), tt.want.item.GetName())
			assert.InDelta(t, tt.want.item.GetPrice(), item.GetPrice(), 0.0001)
		})
	}
}

func TestEntity_CreateCatalogItem(t *testing.T) {
	t.Parallel()
	catalogID := uuid.New().String()
	patterns := []struct {
		name string
		arg  struct {
			name  string
			price float64
		}
		want struct {
			item *CatalogItem
			err  error
		}
	}{
		{
			name: "Success",
			arg: struct {
				name  string
				price float64
			}{
				name:  "item",
				price: 100,
			},
			want: struct {
				item *CatalogItem
				err  error
			}{
				item: &CatalogItem{
					CatalogItem: catalog.CatalogItem{
						Id:    catalogID,
						Name:  "item",
						Price: 100,
					},
				},
				err: nil,
			},
		},
		{
			name: "Fail: name is empty",
			arg: struct {
				name  string
				price float64
			}{
				name:  "",
				price: 100,
			},
			want: struct {
				item *CatalogItem
				err  error
			}{
				item: nil,
				err:  errors.New("name is required"),
			},
		},
		{
			name: "Fail: price is less than 0",
			arg: struct {
				name  string
				price float64
			}{
				name:  "item",
				price: -1,
			},
			want: struct {
				item *CatalogItem
				err  error
			}{
				item: nil,
				err:  errors.New("price must be greater than 0"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			item, err := CreateCatalogItem(tt.arg.name, tt.arg.price)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("CreateCatalogItem() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateCatalogItem() error = %v, wantErr %v", err, tt.want.err)
			}
			if tt.want.err != nil {
				return
			}
			assert.Equal(t, item.GetName(), tt.want.item.GetName())
			assert.InDelta(t, tt.want.item.GetPrice(), item.GetPrice(), 0.0001)
		})
	}
}

func Test_CatalogItem_Proto(t *testing.T) {
	item := CatalogItem{
		CatalogItem: catalog.CatalogItem{
			Id:    "1",
			Name:  "item",
			Price: 100,
		},
	}
	expect := &catalog.CatalogItem{
		Id:    "1",
		Name:  "item",
		Price: 100,
	}
	assert.Equal(t, expect, item.Proto())
}

func Test_CatalogItems_Proto(t *testing.T) {
	items := CatalogItems{
		&CatalogItem{
			catalog.CatalogItem{
				Id:    "1",
				Name:  "item1",
				Price: 100,
			},
		},
		&CatalogItem{
			catalog.CatalogItem{
				Id:    "2",
				Name:  "item2",
				Price: 200,
			},
		},
	}
	expect := []*catalog.CatalogItem{
		{
			Id:    "1",
			Name:  "item1",
			Price: 100,
		},
		{
			Id:    "2",
			Name:  "item2",
			Price: 200,
		},
	}
	assert.Equal(t, expect, items.Proto())
}
