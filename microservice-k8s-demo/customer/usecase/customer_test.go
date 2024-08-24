package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/repository/mock"
)

func TestUseCase_GetCustomer(t *testing.T) {
	t.Parallel()

	customerID := uuid.New().String()

	customer := &entity.Customer{
		ID:      customerID,
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Street:  "123 Maple Street",
		City:    "Springfield",
		Country: "USA",
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		want struct {
			customer *entity.Customer
			err      error
		}
	}{
		{
			name: "success",
			setup: func(cr *mock.MockCustomerRepository) {
				cr.EXPECT().Get(gomock.Any(), customerID).Return(customer, nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  customerID,
			},
			want: struct {
				customer *entity.Customer
				err      error
			}{
				customer: customer,
				err:      nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := mock.NewMockCustomerRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr)
			}

			cuc := NewCustomerUsecase(cr)

			getCustomer, err := cuc.GetCustomer(tt.arg.ctx, tt.arg.id)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetCustomer() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetCustomer() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getCustomer, tt.want.customer) {
				t.Errorf("GetCustomer() got = %v, want %v", getCustomer, tt.want.customer)
			}
		})
	}
}

func TestUseCase_ListCustomers(t *testing.T) {
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
			m *mock.MockCustomerRepository,
		)
		arg struct {
			ctx context.Context
		}
		want struct {
			customers []entity.Customer
			err       error
		}
	}{
		{
			name: "success",
			setup: func(cr *mock.MockCustomerRepository) {
				cr.EXPECT().List(gomock.Any()).Return(customers, nil)
			},
			arg: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			want: struct {
				customers []entity.Customer
				err       error
			}{
				customers: customers,
				err:       nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := mock.NewMockCustomerRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr)
			}

			cuc := NewCustomerUsecase(cr)

			getCustomers, err := cuc.ListCustomers(tt.arg.ctx)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("ListCustomers() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("ListCustomers() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getCustomers, tt.want.customers) {
				t.Errorf("ListCustomers() got = %v, want %v", getCustomers, tt.want.customers)
			}
		})
	}
}

func TestUseCase_CreateCustomer(t *testing.T) { //nolint: gocognit // This is a test function
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerRepository,
		)
		arg struct {
			ctx    context.Context
			params *CreateCustomerParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(cr *mock.MockCustomerRepository) {
				cr.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, customer entity.Customer) {
					if customer.Name != "John Doe" {
						t.Errorf("unexpected Name: got %v, want %v", customer.Name, "John Doe")
					}
					if customer.Email != "john.doe@example.com" {
						t.Errorf("unexpected Email: got %v, want %v", customer.Email, "john.doe@example.com")
					}
					if customer.Street != "123 Maple Street" {
						t.Errorf("unexpected Street: got %v, want %v", customer.Street, "123 Maple Street")
					}
					if customer.City != "Springfield" {
						t.Errorf("unexpected City: got %v, want %v", customer.City, "Springfield")
					}
					if customer.Country != "USA" {
						t.Errorf("unexpected Country: got %v, want %v", customer.Country, "USA")
					}
				}).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params *CreateCustomerParams
			}{
				ctx: context.Background(),
				params: &CreateCustomerParams{
					Name:    "John Doe",
					Email:   "john.doe@example.com",
					Street:  "123 Maple Street",
					City:    "Springfield",
					Country: "USA",
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
			cr := mock.NewMockCustomerRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr)
			}

			cuc := NewCustomerUsecase(cr)

			err := cuc.CreateCustomer(tt.arg.ctx, tt.arg.params)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUseCase_UpdateCustomer(t *testing.T) { //nolint: gocognit // This is a test function
	t.Parallel()

	customerID := uuid.New().String()

	customer := &entity.Customer{
		ID:      customerID,
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Street:  "123 Maple Street",
		City:    "Springfield",
		Country: "USA",
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerRepository,
		)
		arg struct {
			ctx    context.Context
			params *UpdateCustomerParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(cr *mock.MockCustomerRepository) {
				cr.EXPECT().Get(
					gomock.Any(),
					customerID,
				).Return(customer, nil)
				cr.EXPECT().Update(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, customer entity.Customer) {
					if customer.Name != "John Doe" {
						t.Errorf("unexpected Name: got %v, want %v", customer.Name, "John Doe")
					}
					if customer.Email != "john.doe@example.com" {
						t.Errorf("unexpected Email: got %v, want %v", customer.Email, "john.doe@example.com")
					}
					if customer.Street != "123 Maple Street" {
						t.Errorf("unexpected Street: got %v, want %v", customer.Street, "123 Maple Street")
					}
					if customer.City != "Springfield" {
						t.Errorf("unexpected City: got %v, want %v", customer.City, "Springfield")
					}
					if customer.Country != "USA" {
						t.Errorf("unexpected Country: got %v, want %v", customer.Country, "USA")
					}
				}).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params *UpdateCustomerParams
			}{
				ctx: context.Background(),
				params: &UpdateCustomerParams{
					ID:      customerID,
					Name:    "John Doe",
					Email:   "john.doe@example.com",
					Street:  "123 Maple Street",
					City:    "Springfield",
					Country: "USA",
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
			cr := mock.NewMockCustomerRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr)
			}

			cuc := NewCustomerUsecase(cr)

			err := cuc.UpdateCustomer(tt.arg.ctx, tt.arg.params)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUsaCase_DeleteCustomer(t *testing.T) {
	t.Parallel()

	customerID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCustomerRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(cr *mock.MockCustomerRepository) {
				cr.EXPECT().Delete(gomock.Any(), customerID).Return(nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: context.Background(),
				id:  customerID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cr := mock.NewMockCustomerRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr)
			}

			cuc := NewCustomerUsecase(cr)

			err := cuc.DeleteCustomer(tt.arg.ctx, tt.arg.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
