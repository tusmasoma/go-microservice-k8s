package service

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
// 	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository/mock"
// )

// func TestOrderService_CreateOrder(t *testing.T) {
// 	t.Parallel()

// 	orderID := uuid.New().String()

// 	customerID := uuid.New().String()
// 	customer := entity.Customer{
// 		ID: customerID,
// 	}

// 	catalogItemID := uuid.New().String()
// 	catalogItem := entity.CatalogItem{
// 		ID:    catalogItemID,
// 		Name:  "item1",
// 		Price: 1000,
// 	}

// 	orderLines := []entity.OrderLine{
// 		{
// 			Count:       1,
// 			CatalogItem: catalogItem,
// 		},
// 	}

// 	date := time.Now()
// 	order := entity.Order{
// 		ID:         orderID,
// 		Customer:   customer,
// 		OrderDate:  date,
// 		OrderLines: orderLines,
// 	}

// 	patterns := []struct {
// 		name  string
// 		setup func(
// 			m *mock.MockOrderRepository,
// 			m1 *mock.MockOrderLineRepository,
// 			m2 *mock.MockTransactionRepository,
// 		)
// 		arg struct {
// 			ctx   context.Context
// 			order entity.Order
// 		}
// 		wantErr error
// 	}{
// 		{
// 			name: "success",
// 			setup: func(
// 				m *mock.MockOrderRepository,
// 				m1 *mock.MockOrderLineRepository,
// 				m2 *mock.MockTransactionRepository,
// 			) {
// 				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
// 					return fn(ctx)
// 				})
// 				m.EXPECT().Create(
// 					gomock.Any(),
// 					entity.OrderModel{
// 						ID:         orderID,
// 						CustomerID: customerID,
// 						OrderDate:  date,
// 					},
// 				).Return(nil)
// 				m1.EXPECT().BatchCreate(
// 					gomock.Any(),
// 					[]entity.OrderLineModel{
// 						{
// 							OrderID:       orderID,
// 							CatalogItemID: catalogItemID,
// 							Count:         1,
// 						},
// 					},
// 				).Return(nil)
// 			},
// 			arg: struct {
// 				ctx   context.Context
// 				order entity.Order
// 			}{
// 				ctx:   context.Background(),
// 				order: order,
// 			},
// 			wantErr: nil,
// 		},
// 	}
// 	for _, tt := range patterns {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)
// 			or := mock.NewMockOrderRepository(ctrl)
// 			olr := mock.NewMockOrderLineRepository(ctrl)
// 			tr := mock.NewMockTransactionRepository(ctrl)

// 			if tt.setup != nil {
// 				tt.setup(or, olr, tr)
// 			}

// 			service := NewOrderService(or, olr, tr)

// 			err := service.CreateOrder(tt.arg.ctx, tt.arg.order)

// 			if (err != nil) != (tt.wantErr != nil) {
// 				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
// 			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
// 				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestOrderService_DeleteOrder(t *testing.T) {
// 	t.Parallel()

// 	orderID := uuid.New().String()

// 	patterns := []struct {
// 		name  string
// 		setup func(
// 			m *mock.MockOrderRepository,
// 			m1 *mock.MockOrderLineRepository,
// 			m2 *mock.MockTransactionRepository,
// 		)
// 		arg struct {
// 			ctx context.Context
// 			id  string
// 		}
// 		wantErr error
// 	}{
// 		{
// 			name: "success",
// 			setup: func(
// 				m *mock.MockOrderRepository,
// 				m1 *mock.MockOrderLineRepository,
// 				m2 *mock.MockTransactionRepository,
// 			) {
// 				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
// 					return fn(ctx)
// 				})
// 				m1.EXPECT().BatchDelete(
// 					gomock.Any(),
// 					orderID,
// 				).Return(nil)
// 				m.EXPECT().Delete(
// 					gomock.Any(),
// 					orderID,
// 				).Return(nil)
// 			},
// 			arg: struct {
// 				ctx context.Context
// 				id  string
// 			}{
// 				ctx: context.Background(),
// 				id:  orderID,
// 			},
// 			wantErr: nil,
// 		},
// 	}
// 	for _, tt := range patterns {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)
// 			or := mock.NewMockOrderRepository(ctrl)
// 			olr := mock.NewMockOrderLineRepository(ctrl)
// 			tr := mock.NewMockTransactionRepository(ctrl)

// 			if tt.setup != nil {
// 				tt.setup(or, olr, tr)
// 			}

// 			service := NewOrderService(or, olr, tr)

// 			err := service.DeleteOrder(tt.arg.ctx, tt.arg.id)

// 			if (err != nil) != (tt.wantErr != nil) {
// 				t.Errorf("DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
// 			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
// 				t.Errorf("DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
