package entity

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
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
			name: "success: id is not empty",
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
					ID:    catalogID,
					Name:  "item",
					Price: 100,
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
				item: &CatalogItem{
					ID:    uuid.New().String(),
					Name:  "item",
					Price: 100,
				},
				err: nil,
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

			if d := cmp.Diff(item, tt.want.item, cmpopts.IgnoreFields(CatalogItem{}, "ID")); len(d) != 0 {
				t.Errorf("NewCatalogItem() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
