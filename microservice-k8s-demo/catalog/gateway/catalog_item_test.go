package gateway

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/proto"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/usecase/mock"
)

const bufSize = 1024 * 1024

func setupTestServer(t *testing.T, setup func(m *mock.MockCatalogItemUseCase)) (pb.CatalogServiceClient, func()) {
	t.Helper()

	ctrl := gomock.NewController(t)
	cuc := mock.NewMockCatalogItemUseCase(ctrl)

	if setup != nil {
		setup(cuc)
	}

	handler := NewCatalogItemHandler(cuc)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, handler)

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

	client := pb.NewCatalogServiceClient(conn)

	cleanup := func() {
		conn.Close()
		s.Stop()
	}

	return client, cleanup
}

func TestHandler_ListCatalogItemsByName(t *testing.T) {
	t.Parallel()

	items := []entity.CatalogItem{
		{
			ID:    uuid.New().String(),
			Name:  "item1",
			Price: float64(100),
		},
		{
			ID:    uuid.New().String(),
			Name:  "item2",
			Price: float64(200),
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemUseCase,
		)
		request    *pb.ListCatalogItemsByNameRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(cuc *mock.MockCatalogItemUseCase) {
				cuc.EXPECT().ListCatalogItemsByName(
					gomock.Any(),
					"item",
				).Return(items, nil)
			},
			request: &pb.ListCatalogItemsByNameRequest{
				Name: "item",
			},
			wantStatus: codes.OK,
		},
		{
			name:       "Fail: invalid request of name is empty",
			request:    &pb.ListCatalogItemsByNameRequest{Name: ""},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			resp, err := client.ListCatalogItemsByName(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if len(resp.GetItems()) != len(items) {
					t.Fatalf("handler returned wrong item data")
				}
			}
		})
	}
}

func TestHandler_ListCatalogItems(t *testing.T) {
	t.Parallel()

	items := []entity.CatalogItem{
		{
			ID:    uuid.New().String(),
			Name:  "item1",
			Price: float64(100),
		},
		{
			ID:    uuid.New().String(),
			Name:  "item2",
			Price: float64(200),
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemUseCase,
		)
		request    *pb.ListCatalogItemsRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCatalogItemUseCase) {
				tuc.EXPECT().ListCatalogItems(
					gomock.Any(),
				).Return(items, nil)
			},
			request:    &pb.ListCatalogItemsRequest{},
			wantStatus: codes.OK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			resp, err := client.ListCatalogItems(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if len(resp.GetItems()) != len(items) {
					t.Fatalf("handler returned wrong item data")
				}
			}
		})
	}
}

func TestHandler_CreateCatalogItem(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemUseCase,
		)
		request    *pb.CreateCatalogItemRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCatalogItemUseCase) {
				tuc.EXPECT().CreateCatalogItem(
					gomock.Any(),
					"item1",
					float64(100),
				).Return(nil)
			},
			request: &pb.CreateCatalogItemRequest{
				Name:  "item1",
				Price: float64(100),
			},
			wantStatus: codes.OK,
		},
		{
			name: "Fail: invalid request of name is empty",
			request: &pb.CreateCatalogItemRequest{
				Name:  "",
				Price: float64(100),
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of  price is less than 0",
			request: &pb.CreateCatalogItemRequest{
				Name:  "item1",
				Price: float64(-1),
			},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			req, err := client.CreateCatalogItem(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if req == nil {
					t.Fatalf("handler returned wrong item data")
				}
			}
		})
	}
}

func TestHandler_UpdateCatalogItem(t *testing.T) {
	t.Parallel()

	itemID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemUseCase,
		)
		request    *pb.UpdateCatalogItemRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCatalogItemUseCase) {
				tuc.EXPECT().UpdateCatalogItem(
					gomock.Any(),
					itemID,
					"updated name",
					float64(100),
				).Return(nil)
			},
			request: &pb.UpdateCatalogItemRequest{
				Id:    itemID,
				Name:  "updated name",
				Price: float64(100),
			},
			wantStatus: codes.OK,
		},
		{
			name: "Fail: invalid request of id is empty",
			request: &pb.UpdateCatalogItemRequest{
				Id:    "",
				Name:  "updated name",
				Price: float64(100),
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of name is empty",
			request: &pb.UpdateCatalogItemRequest{
				Id:    itemID,
				Name:  "",
				Price: float64(100),
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of price is less than 0",
			request: &pb.UpdateCatalogItemRequest{
				Id:    itemID,
				Name:  "updated name",
				Price: float64(-1),
			},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			req, err := client.UpdateCatalogItem(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if req == nil {
					t.Fatalf("handler returned wrong item data")
				}
			}
		})
	}
}

func TestHandler_DeleteCatalogItem(t *testing.T) {
	t.Parallel()

	itemID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCatalogItemUseCase,
		)
		request    *pb.DeleteCatalogItemRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCatalogItemUseCase) {
				tuc.EXPECT().DeleteCatalogItem(
					gomock.Any(),
					itemID,
				).Return(nil)
			},
			request:    &pb.DeleteCatalogItemRequest{Id: itemID},
			wantStatus: codes.OK,
		},
		{
			name:       "Fail: invalid request of id is empty",
			request:    &pb.DeleteCatalogItemRequest{Id: ""},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			req, err := client.DeleteCatalogItem(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if req == nil {
					t.Fatalf("handler returned wrong item data")
				}
			}
		})
	}
}
