package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/repository/mock"
)

func TestUseCase_GetCatalogItem(t *testing.T) {
	t.Parallel()

	itemID := uuid.New().String()

	item := &entity.CatalogItem{
		ID:    itemID,
		Name:  "item",
		Price: 100,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		want struct {
			item *entity.CatalogItem
			err  error
		}
	}{
		{
			name: "success",
			setup: func(tr *mock.MockCatalogItemRepository) {
				tr.EXPECT().Get(gomock.Any(), itemID).Return(item, nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  itemID,
			},
			want: struct {
				item *entity.CatalogItem
				err  error
			}{
				item: item,
				err:  nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockCatalogItemRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewCatalogItemUseCase(tr)

			getCatalogItem, err := tuc.GetCatalogItem(tt.arg.ctx, tt.arg.id)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetCatalogItem() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetCatalogItem() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getCatalogItem, tt.want.item) {
				t.Errorf("GetCatalogItem() got = %v, want %v", getCatalogItem, tt.want.item)
			}
		})
	}
}

func TestUseCase_ListCatalogItems(t *testing.T) {
	t.Parallel()

	items := []entity.CatalogItem{
		{
			ID:    uuid.New().String(),
			Name:  "item1",
			Price: 100,
		},
		{
			ID:    uuid.New().String(),
			Name:  "item2",
			Price: 200,
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemRepository,
		)
		arg struct {
			ctx context.Context
		}
		want struct {
			items []entity.CatalogItem
			err   error
		}
	}{
		{
			name: "success",
			setup: func(tr *mock.MockCatalogItemRepository) {
				tr.EXPECT().List(gomock.Any()).Return(items, nil)
			},
			arg: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			want: struct {
				items []entity.CatalogItem
				err   error
			}{
				items: items,
				err:   nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockCatalogItemRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewCatalogItemUseCase(tr)

			getCatalogItems, err := tuc.ListCatalogItems(tt.arg.ctx)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("ListCatalogItems() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("ListCatalogItems() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getCatalogItems, tt.want.items) {
				t.Errorf("ListCatalogItems() got = %v, want %v", getCatalogItems, tt.want.items)
			}
		})
	}
}

func TestUseCase_ListCatalogItemsByName(t *testing.T) {
	t.Parallel()

	item1 := entity.CatalogItem{
		ID:    uuid.New().String(),
		Name:  "item1",
		Price: 100,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemRepository,
		)
		arg struct {
			ctx  context.Context
			name string
		}
		want struct {
			items []entity.CatalogItem
			err   error
		}
	}{
		{
			name: "success",
			setup: func(tr *mock.MockCatalogItemRepository) {
				tr.EXPECT().ListByName(gomock.Any(), "item1").Return(
					[]entity.CatalogItem{item1},
					nil,
				)
			},
			arg: struct {
				ctx  context.Context
				name string
			}{
				ctx:  context.Background(),
				name: "item1",
			},
			want: struct {
				items []entity.CatalogItem
				err   error
			}{
				items: []entity.CatalogItem{item1},
				err:   nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockCatalogItemRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewCatalogItemUseCase(tr)

			getCatalogItems, err := tuc.ListCatalogItemsByName(tt.arg.ctx, tt.arg.name)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("ListCatalogItemsByName() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("ListCatalogItemsByName() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getCatalogItems, tt.want.items) {
				t.Errorf("ListCatalogItemsByName() got = %v, want %v", getCatalogItems, tt.want.items)
			}
		})
	}
}

func TestUseCase_ListCatalogItemsByIDs(t *testing.T) {
	t.Parallel()

	item1 := entity.CatalogItem{
		ID:    uuid.New().String(),
		Name:  "item1",
		Price: 100,
	}
	item2 := entity.CatalogItem{
		ID:    uuid.New().String(),
		Name:  "item2",
		Price: 200,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemRepository,
		)
		arg struct {
			ctx context.Context
			ids []string
		}
		want struct {
			items []entity.CatalogItem
			err   error
		}
	}{
		{
			name: "success",
			setup: func(tr *mock.MockCatalogItemRepository) {
				tr.EXPECT().ListByIDs(gomock.Any(), []string{item1.ID, item2.ID}).Return(
					[]entity.CatalogItem{item1, item2},
					nil,
				)
			},
			arg: struct {
				ctx context.Context
				ids []string
			}{
				ctx: context.Background(),
				ids: []string{item1.ID, item2.ID},
			},
			want: struct {
				items []entity.CatalogItem
				err   error
			}{
				items: []entity.CatalogItem{item1, item2},
				err:   nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockCatalogItemRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewCatalogItemUseCase(tr)

			getCatalogItems, err := tuc.ListCatalogItemsByIDs(tt.arg.ctx, tt.arg.ids)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("ListCatalogItemsByIDs() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("ListCatalogItemsByName() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getCatalogItems, tt.want.items) {
				t.Errorf("ListCatalogItemsByIDs() got = %v, want %v", getCatalogItems, tt.want.items)
			}
		})
	}
}

func TestUseCase_CreateCatalogItem(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemRepository,
		)
		arg struct {
			ctx   context.Context
			name  string
			price float64
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(tr *mock.MockCatalogItemRepository) {
				tr.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, item entity.CatalogItem) {
					if item.Name != "item" {
						t.Errorf("unexpected Name: got %v, want %v", item.Name, "item")
					}
					if item.Price != 100 {
						t.Errorf("unexpected Price: got %v, want %v", item.Price, 100)
					}
				}).Return(nil)
			},
			arg: struct {
				ctx   context.Context
				name  string
				price float64
			}{
				ctx:   context.Background(),
				name:  "item",
				price: 100,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockCatalogItemRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewCatalogItemUseCase(tr)

			err := tuc.CreateCatalogItem(tt.arg.ctx, tt.arg.name, tt.arg.price)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUseCase_UpdateCatalogItem(t *testing.T) {
	t.Parallel()

	itemID := uuid.New().String()

	item := &entity.CatalogItem{
		ID:    itemID,
		Name:  "item",
		Price: 100,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemRepository,
		)
		arg struct {
			ctx   context.Context
			id    string
			name  string
			price float64
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(tr *mock.MockCatalogItemRepository) {
				tr.EXPECT().Get(
					gomock.Any(),
					itemID,
				).Return(item, nil)
				tr.EXPECT().Update(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, item entity.CatalogItem) {
					if item.Name != "updated item" {
						t.Errorf("unexpected Name: got %v, want %v", item.Name, "updated item")
					}
					if item.Price != 200 {
						t.Errorf("unexpected Price: got %v, want %v", item.Price, 200)
					}
				}).Return(nil)
			},
			arg: struct {
				ctx   context.Context
				id    string
				name  string
				price float64
			}{
				ctx:   context.Background(),
				id:    itemID,
				name:  "updated item",
				price: 200,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockCatalogItemRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewCatalogItemUseCase(tr)

			err := tuc.UpdateCatalogItem(tt.arg.ctx, tt.arg.id, tt.arg.name, tt.arg.price)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUsaCase_DeleteCatalogItem(t *testing.T) {
	t.Parallel()

	itemID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(tr *mock.MockCatalogItemRepository) {
				tr.EXPECT().Delete(gomock.Any(), itemID).Return(nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  itemID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockCatalogItemRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewCatalogItemUseCase(tr)

			err := tuc.DeleteCatalogItem(tt.arg.ctx, tt.arg.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
