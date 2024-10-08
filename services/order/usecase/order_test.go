package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tusmasoma/go-microservice-k8s/services/order/entity"
	repo_mock "github.com/tusmasoma/go-microservice-k8s/services/order/repository/mock"
)

func TestOrderUseCase_GetOrderCreationResources(t *testing.T) {
	t.Parallel()

	customers := []entity.Customer{
		{
			ID:   uuid.New().String(),
			Name: "customer1",
		},
	}

	items := []entity.CatalogItem{
		{
			ID:    uuid.New().String(),
			Name:  "item1",
			Price: 1000,
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *repo_mock.MockCustomerRepository,
			m1 *repo_mock.MockCatalogItemRepository,
			m2 *repo_mock.MockOrderRepository,
		)
		arg struct {
			ctx context.Context
		}
		want struct {
			customers []entity.Customer
			items     []entity.CatalogItem
			err       error
		}
	}{
		{
			name: "success",
			setup: func(
				cr *repo_mock.MockCustomerRepository,
				cir *repo_mock.MockCatalogItemRepository,
				or *repo_mock.MockOrderRepository,
			) {
				cr.EXPECT().List(gomock.Any()).Return(customers, nil)
				cir.EXPECT().List(gomock.Any()).Return(items, nil)
			},
			arg: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			want: struct {
				customers []entity.Customer
				items     []entity.CatalogItem
				err       error
			}{
				customers: customers,
				items:     items,
				err:       nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := repo_mock.NewMockCustomerRepository(ctrl)
			cir := repo_mock.NewMockCatalogItemRepository(ctrl)
			or := repo_mock.NewMockOrderRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr, cir, or)
			}

			ouc := NewOrderUseCase(cr, cir, or)

			gotCustomers, gotItems, err := ouc.GetOrderCreationResources(tt.arg.ctx)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetOrderCreationResources() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetOrderCreationResources() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(gotCustomers, tt.want.customers) {
				t.Errorf("GetOrderCreationResources() got = %v, want %v", gotCustomers, tt.want.customers)
			}
			if !reflect.DeepEqual(gotItems, tt.want.items) {
				t.Errorf("GetOrderCreationResources() got = %v, want %v", gotItems, tt.want.items)
			}
		})
	}
}

func TestOrderUseCase_GetOrder(t *testing.T) {
	t.Parallel()

	orderID := uuid.New().String()
	customerID := uuid.New().String()
	catalogItemID := uuid.New().String()
	orderDate := time.Now()

	order := entity.Order{
		ID:         orderID,
		CustomerID: customerID,
		OrderDate:  &orderDate,
		OrderLines: []*entity.OrderLine{
			{
				Count:         1,
				CatalogItemID: catalogItemID,
			},
		},
	}

	customer := entity.Customer{
		ID:   customerID,
		Name: "customer1",
	}

	item := entity.CatalogItem{
		ID:    catalogItemID,
		Name:  "item1",
		Price: 1000,
	}

	patterns := []struct {
		name  string
		setup func(
			m *repo_mock.MockCustomerRepository,
			m1 *repo_mock.MockCatalogItemRepository,
			m2 *repo_mock.MockOrderRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		want struct {
			orderDetails *OrderDetails
			err          error
		}
	}{
		{
			name: "success",
			setup: func(
				cr *repo_mock.MockCustomerRepository,
				cir *repo_mock.MockCatalogItemRepository,
				or *repo_mock.MockOrderRepository,
			) {
				or.EXPECT().Get(gomock.Any(), orderID).Return(
					&order,
					nil,
				)
				cr.EXPECT().Get(gomock.Any(), customerID).Return(
					&customer, nil)
				cir.EXPECT().ListByIDs(gomock.Any(), []string{catalogItemID}).Return(
					[]entity.CatalogItem{item}, nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  orderID,
			},
			want: struct {
				orderDetails *OrderDetails
				err          error
			}{
				orderDetails: &OrderDetails{
					Order:    &order,
					Customer: &customer,
					OrderLines: []*OrderLineDetails{
						{
							Count:       1,
							CatalogItem: &item,
						},
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := repo_mock.NewMockCustomerRepository(ctrl)
			cir := repo_mock.NewMockCatalogItemRepository(ctrl)
			or := repo_mock.NewMockOrderRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr, cir, or)
			}

			ouc := NewOrderUseCase(cr, cir, or)

			gotOrderDetails, err := ouc.GetOrder(tt.arg.ctx, tt.arg.id)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetOrder() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetOrder() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(gotOrderDetails, tt.want.orderDetails) {
				t.Errorf("GetOrder() got = %v, want %v", gotOrderDetails, tt.want.orderDetails)
			}
		})
	}
}

func TestOrderUseCase_ListOrder(t *testing.T) {
	t.Parallel()

	orderID := uuid.New().String()
	customerID := uuid.New().String()
	catalogItemID := uuid.New().String()
	orderDate := time.Now()

	orders := []*entity.Order{
		{
			ID:         orderID,
			CustomerID: customerID,
			OrderDate:  &orderDate,
			OrderLines: []*entity.OrderLine{
				{
					Count:         1,
					CatalogItemID: catalogItemID,
				},
			},
		},
	}

	customer := entity.Customer{
		ID:   customerID,
		Name: "customer1",
	}

	item := entity.CatalogItem{
		ID:    catalogItemID,
		Name:  "item1",
		Price: 1000,
	}

	patterns := []struct {
		name  string
		setup func(
			m *repo_mock.MockCustomerRepository,
			m1 *repo_mock.MockCatalogItemRepository,
			m2 *repo_mock.MockOrderRepository,
		)
		arg struct {
			ctx context.Context
		}
		want struct {
			orderDetails []*OrderDetails
			err          error
		}
	}{
		{
			name: "success",
			setup: func(
				cr *repo_mock.MockCustomerRepository,
				cir *repo_mock.MockCatalogItemRepository,
				or *repo_mock.MockOrderRepository,
			) {
				or.EXPECT().List(gomock.Any()).Return(
					orders,
					nil,
				)
				cr.EXPECT().Get(gomock.Any(), customerID).Return(
					&customer, nil)
				cir.EXPECT().ListByIDs(gomock.Any(), []string{catalogItemID}).Return(
					[]entity.CatalogItem{item}, nil)
			},
			arg: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			want: struct {
				orderDetails []*OrderDetails
				err          error
			}{
				orderDetails: []*OrderDetails{
					{
						Order:    orders[0],
						Customer: &customer,
						OrderLines: []*OrderLineDetails{
							{
								Count:       1,
								CatalogItem: &item,
							},
						},
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := repo_mock.NewMockCustomerRepository(ctrl)
			cir := repo_mock.NewMockCatalogItemRepository(ctrl)
			or := repo_mock.NewMockOrderRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr, cir, or)
			}

			ouc := NewOrderUseCase(cr, cir, or)

			gotOrderDetails, err := ouc.ListOrders(tt.arg.ctx)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("ListOrder() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("ListOrder() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(gotOrderDetails, tt.want.orderDetails) {
				t.Errorf("ListOrder() got = %v, want %v", gotOrderDetails, tt.want.orderDetails)
			}
		})
	}
}

func TestOrderUseCase_CreateOrder(t *testing.T) {
	t.Parallel()

	customerID := uuid.New().String()
	catalogItemID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *repo_mock.MockCustomerRepository,
			m1 *repo_mock.MockCatalogItemRepository,
			m2 *repo_mock.MockOrderRepository,
		)
		arg struct {
			ctx    context.Context
			params *CreateOrderParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				cr *repo_mock.MockCustomerRepository,
				cir *repo_mock.MockCatalogItemRepository,
				or *repo_mock.MockOrderRepository,
			) {
				or.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, order entity.Order) {
					if order.CustomerID != customerID {
						t.Errorf("unexpected customerID: got %v, want %v", order.CustomerID, customerID)
					}
					if order.OrderLines[0].CatalogItemID != catalogItemID {
						t.Errorf("unexpected catalogItemID: got %v, want %v", order.OrderLines[0].CatalogItemID, catalogItemID)
					}
				}).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params *CreateOrderParams
			}{
				ctx: context.Background(),
				params: &CreateOrderParams{
					CustomerID: customerID,
					OrderLine: []struct {
						CatalogItemID string
						Count         int
					}{
						{
							CatalogItemID: catalogItemID,
							Count:         1,
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := repo_mock.NewMockCustomerRepository(ctrl)
			cir := repo_mock.NewMockCatalogItemRepository(ctrl)
			or := repo_mock.NewMockOrderRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr, cir, or)
			}

			ouc := NewOrderUseCase(cr, cir, or)

			err := ouc.CreateOrder(tt.arg.ctx, tt.arg.params)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderUseCase_DeleteOrder(t *testing.T) {
	t.Parallel()

	orderID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *repo_mock.MockCustomerRepository,
			m1 *repo_mock.MockCatalogItemRepository,
			m2 *repo_mock.MockOrderRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				cr *repo_mock.MockCustomerRepository,
				cir *repo_mock.MockCatalogItemRepository,
				or *repo_mock.MockOrderRepository,
			) {
				or.EXPECT().Delete(gomock.Any(), orderID).Return(nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  orderID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := repo_mock.NewMockCustomerRepository(ctrl)
			cir := repo_mock.NewMockCatalogItemRepository(ctrl)
			or := repo_mock.NewMockOrderRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr, cir, or)
			}

			ouc := NewOrderUseCase(cr, cir, or)

			err := ouc.DeleteOrder(tt.arg.ctx, tt.arg.id)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
