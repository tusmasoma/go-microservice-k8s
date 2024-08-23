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

	"github.com/tusmasoma/microservice-k8s-demo/customer/entity"
	"github.com/tusmasoma/microservice-k8s-demo/customer/usecase"

	pb "github.com/tusmasoma/microservice-k8s-demo/customer/proto"

	"github.com/tusmasoma/microservice-k8s-demo/customer/usecase/mock"
)

const bufSize = 1024 * 1024

func setupTestServer(t *testing.T, setup func(m *mock.MockCustomerUseCase)) (pb.CustomerServiceClient, func()) {
	t.Helper()

	ctrl := gomock.NewController(t)
	cuc := mock.NewMockCustomerUseCase(ctrl)

	if setup != nil {
		setup(cuc)
	}

	handler := NewCustomerHandler(cuc)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCustomerServiceServer(s, handler)

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

	client := pb.NewCustomerServiceClient(conn)

	cleanup := func() {
		conn.Close()
		s.Stop()
	}

	return client, cleanup
}

func TestHandler_ListCustomers(t *testing.T) {
	t.Parallel()

	customers := []entity.Customer{
		{
			ID:      uuid.New().String(),
			Name:    "John Doe",
			Email:   "john.doe@example.com",
			Street:  "123 Maple Street",
			City:    "Springfield",
			Country: "USA",
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerUseCase,
		)
		request    *pb.ListCustomersRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCustomerUseCase) {
				tuc.EXPECT().ListCustomers(
					gomock.Any(),
				).Return(customers, nil)
			},
			request:    &pb.ListCustomersRequest{},
			wantStatus: codes.OK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			resp, err := client.ListCustomers(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if len(resp.GetCustomers()) != len(customers) {
					t.Fatalf("handler returned wrong item data")
				}
			}
		})
	}
}

func TestHandler_CreateCustomer(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerUseCase,
		)
		request    *pb.CreateCustomerRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCustomerUseCase) {
				tuc.EXPECT().CreateCustomer(
					gomock.Any(),
					&usecase.CreateCustomerParams{
						Name:    "John Doe",
						Email:   "john.doe@example.com",
						Street:  "123 Maple Street",
						City:    "Springfield",
						Country: "USA",
					},
				).Return(nil)
			},
			request: &pb.CreateCustomerRequest{
				Name:    "John Doe",
				Email:   "john.doe@example.com",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.OK,
		},
		{
			name: "Fail: invalid request of name is empty",
			request: &pb.CreateCustomerRequest{
				Name:    "",
				Email:   "john.doe@example.com",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of  email is empty",
			request: &pb.CreateCustomerRequest{
				Name:    "John Doe",
				Email:   "",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of street is empty",
			request: &pb.CreateCustomerRequest{
				Name:    "John Doe",
				Email:   "john.doe@example.com",
				Street:  "",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of city is empty",
			request: &pb.CreateCustomerRequest{
				Name:    "John Doe",
				Email:   "john.doe@example.com",
				Street:  "123 Maple Street",
				City:    "",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of country is empty",
			request: &pb.CreateCustomerRequest{
				Name:    "John Doe",
				Email:   "john.doe@example.com",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "",
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

			req, err := client.CreateCustomer(context.Background(), tt.request)
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

func TestHandler_UpdateCustomer(t *testing.T) {
	t.Parallel()

	itemID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerUseCase,
		)
		request    *pb.UpdateCustomerRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCustomerUseCase) {
				tuc.EXPECT().UpdateCustomer(
					gomock.Any(),
					&usecase.UpdateCustomerParams{
						ID:      itemID,
						Name:    "New John Doe",
						Email:   "john.new.doe@example.com",
						Street:  "123 Maple Street",
						City:    "Springfield",
						Country: "USA",
					},
				).Return(nil)
			},
			request: &pb.UpdateCustomerRequest{
				Id:      itemID,
				Name:    "New John Doe",
				Email:   "john.new.doe@example.com",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.OK,
		},
		{
			name: "Fail: invalid request of id is empty",
			request: &pb.UpdateCustomerRequest{
				Id:      "",
				Name:    "New John Doe",
				Email:   "john.new.doe@example.com",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of name is empty",
			request: &pb.UpdateCustomerRequest{
				Id:      itemID,
				Name:    "",
				Email:   "john.new.doe@example.com",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of email is empty",
			request: &pb.UpdateCustomerRequest{
				Id:      itemID,
				Name:    "New John Doe",
				Email:   "",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of street is empty",
			request: &pb.UpdateCustomerRequest{
				Id:      itemID,
				Name:    "New John Doe",
				Email:   "john.new.doe@example.com",
				Street:  "",
				City:    "Springfield",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of city is empty",
			request: &pb.UpdateCustomerRequest{
				Id:      itemID,
				Name:    "New John Doe",
				Email:   "john.new.doe@example.com",
				Street:  "123 Maple Street",
				City:    "",
				Country: "USA",
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of country is empty",
			request: &pb.UpdateCustomerRequest{
				Id:      itemID,
				Name:    "New John Doe",
				Email:   "john.new.doe@example.com",
				Street:  "123 Maple Street",
				City:    "Springfield",
				Country: "",
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

			req, err := client.UpdateCustomer(context.Background(), tt.request)
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

func TestHandler_DeleteCustomer(t *testing.T) {
	t.Parallel()

	itemID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerUseCase,
		)
		request    *pb.DeleteCustomerRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockCustomerUseCase) {
				tuc.EXPECT().DeleteCustomer(
					gomock.Any(),
					itemID,
				).Return(nil)
			},
			request:    &pb.DeleteCustomerRequest{Id: itemID},
			wantStatus: codes.OK,
		},
		{
			name:       "Fail: invalid request of id is empty",
			request:    &pb.DeleteCustomerRequest{Id: ""},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			req, err := client.DeleteCustomer(context.Background(), tt.request)
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
