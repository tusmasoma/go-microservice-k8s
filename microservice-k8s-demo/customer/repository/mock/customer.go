// Code generated by MockGen. DO NOT EDIT.
// Source: customer.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/tusmasoma/microservice-k8s-demo/customer/entity"
)

// MockCustomerRepository is a mock of CustomerRepository interface.
type MockCustomerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCustomerRepositoryMockRecorder
}

// MockCustomerRepositoryMockRecorder is the mock recorder for MockCustomerRepository.
type MockCustomerRepositoryMockRecorder struct {
	mock *MockCustomerRepository
}

// NewMockCustomerRepository creates a new mock instance.
func NewMockCustomerRepository(ctrl *gomock.Controller) *MockCustomerRepository {
	mock := &MockCustomerRepository{ctrl: ctrl}
	mock.recorder = &MockCustomerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomerRepository) EXPECT() *MockCustomerRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCustomerRepository) Create(ctx context.Context, customer entity.Customer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, customer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCustomerRepositoryMockRecorder) Create(ctx, customer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCustomerRepository)(nil).Create), ctx, customer)
}

// Delete mocks base method.
func (m *MockCustomerRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCustomerRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCustomerRepository)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockCustomerRepository) Get(ctx context.Context, id string) (*entity.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCustomerRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCustomerRepository)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockCustomerRepository) List(ctx context.Context) ([]entity.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]entity.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCustomerRepositoryMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCustomerRepository)(nil).List), ctx)
}

// Update mocks base method.
func (m *MockCustomerRepository) Update(ctx context.Context, customer entity.Customer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, customer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCustomerRepositoryMockRecorder) Update(ctx, customer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCustomerRepository)(nil).Update), ctx, customer)
}
