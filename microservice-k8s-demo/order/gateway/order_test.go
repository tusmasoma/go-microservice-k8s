package gateway

import (
	"context"
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/usecase"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/usecase/mock"
)

const bufSize = 1024 * 1024

func setupTestServer(t *testing.T, setup func(m *mock.MockOrderUseCase)) (pb.OrderServiceClient, func()) {
	t.Helper()

	ctrl := gomock.NewController(t)
	cuc := mock.NewMockOrderUseCase(ctrl)

	if setup != nil {
		setup(cuc)
	}

	handler := NewOrderHandler(cuc)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, handler)

	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.Dial("", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { //nolint:staticcheck // ignore deprecation
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	client := pb.NewOrderServiceClient(conn)

	cleanup := func() {
		conn.Close()
		s.Stop()
	}

	return client, cleanup
}

func TestHandler_ListOrders(t *testing.T) {
	t.Parallel()

	orderID := uuid.New().String()
	customerID := uuid.New().String()
	itemID := uuid.New().String()

	date := time.Now()

	customer := entity.Customer{
		ID:      customerID,
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Street:  "123 Maple Street",
		City:    "Springfield",
		Country: "USA",
	}

	item := entity.CatalogItem{
		ID:    itemID,
		Name:  "item",
		Price: 1000,
	}

	order := &entity.Order{
		ID:        orderID,
		Customer:  customer,
		OrderDate: date,
		OrderLines: []entity.OrderLine{
			{
				Count:       1,
				CatalogItem: item,
			},
		},
		TotalPrice: 1000,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockOrderUseCase,
		)
		request    *pb.ListOrdersRequest
		wantStatus codes.Code
		want       []*pb.Order
	}{
		{
			name: "success",
			setup: func(ouc *mock.MockOrderUseCase) {
				ouc.EXPECT().ListOrders(
					gomock.Any(),
				).Return(
					[]*entity.Order{order},
					nil,
				)
			},
			request:    &pb.ListOrdersRequest{},
			wantStatus: codes.OK,
			want: []*pb.Order{
				{
					Id: orderID,
					Customer: &pb.Customer{
						Id:      customerID,
						Name:    "John Doe",
						Email:   "john.doe@example.com",
						Street:  "123 Maple Street",
						City:    "Springfield",
						Country: "USA",
					},
					OrderDate: timestamppb.New(date),
					OrderLines: []*pb.OrderLine{
						{
							Item: &pb.CatalogItem{
								Id:    itemID,
								Name:  "item",
								Price: 1000,
							},
							Count: 1,
						},
					},
					TotalPrice: 1000,
				},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			resp, err := client.ListOrders(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if !reflect.DeepEqual(resp.GetOrders(), tt.want) {
					t.Errorf("handler returned wrong orders: got %v want %v", resp.GetOrders(), tt.want)
				}
			}
		})
	}
}

func TestHandler_GetOrderCreationResources(t *testing.T) {
	t.Parallel()

	customerID := uuid.New().String()
	itemID := uuid.New().String()

	customer := entity.Customer{
		ID:      customerID,
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Street:  "123 Maple Street",
		City:    "Springfield",
		Country: "USA",
	}

	item := entity.CatalogItem{
		ID:    itemID,
		Name:  "item",
		Price: 1000,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockOrderUseCase,
		)
		request    *pb.GetOrderCreationResourcesRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(ouc *mock.MockOrderUseCase) {
				ouc.EXPECT().GetOrderCreationResources(
					gomock.Any(),
				).Return(
					[]entity.Customer{customer},
					[]entity.CatalogItem{item},
					nil,
				)
			},
			request:    &pb.GetOrderCreationResourcesRequest{},
			wantStatus: codes.OK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			resp, err := client.GetOrderCreationResources(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if len(resp.GetCustomers()) != 1 {
					t.Errorf("handler returned wrong number of customers: got %v want %v", len(resp.GetCustomers()), 1)
				}
				if len(resp.GetItems()) != 1 {
					t.Errorf("handler returned wrong number of items: got %v want %v", len(resp.GetItems()), 1)
				}
			}
		})
	}
}

func TestHandler_CreateOrder(t *testing.T) {
	t.Parallel()

	customerID := uuid.New().String()
	itemID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockOrderUseCase,
		)
		request    *pb.CreateOrderRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(ouc *mock.MockOrderUseCase) {
				ouc.EXPECT().CreateOrder(
					gomock.Any(),
					&usecase.CreateOrderParams{
						CustomerID: customerID,
						OrderLine: []struct {
							CatalogItemID string
							Count         int
						}{
							{
								CatalogItemID: itemID,
								Count:         1,
							},
						},
					},
				).Return(nil)
			},
			request: &pb.CreateOrderRequest{
				CustomerId: customerID,
				OrderLines: []*pb.OrderLine{
					{
						Item: &pb.CatalogItem{
							Id: itemID,
						},
						Count: 1,
					},
				},
			},
			wantStatus: codes.OK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			_, err := client.CreateOrder(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}
		})
	}
}

func TestHandler_DeleteOrder(t *testing.T) {
	t.Parallel()

	orderID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockOrderUseCase,
		)
		request    *pb.DeleteOrderRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(ouc *mock.MockOrderUseCase) {
				ouc.EXPECT().DeleteOrder(
					gomock.Any(),
					orderID,
				).Return(nil)
			},
			request: &pb.DeleteOrderRequest{
				OrderId: orderID,
			},
			wantStatus: codes.OK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			_, err := client.DeleteOrder(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}
		})
	}
}
