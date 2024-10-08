// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/tusmasoma/go-microservice-k8s/services/order/entity"
	usecase "github.com/tusmasoma/go-microservice-k8s/services/order/usecase"
)

// MockOrderUseCase is a mock of OrderUseCase interface.
type MockOrderUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUseCaseMockRecorder
}

// MockOrderUseCaseMockRecorder is the mock recorder for MockOrderUseCase.
type MockOrderUseCaseMockRecorder struct {
	mock *MockOrderUseCase
}

// NewMockOrderUseCase creates a new mock instance.
func NewMockOrderUseCase(ctrl *gomock.Controller) *MockOrderUseCase {
	mock := &MockOrderUseCase{ctrl: ctrl}
	mock.recorder = &MockOrderUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUseCase) EXPECT() *MockOrderUseCaseMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrderUseCase) CreateOrder(ctx context.Context, params *usecase.CreateOrderParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderUseCaseMockRecorder) CreateOrder(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderUseCase)(nil).CreateOrder), ctx, params)
}

// DeleteOrder mocks base method.
func (m *MockOrderUseCase) DeleteOrder(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockOrderUseCaseMockRecorder) DeleteOrder(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockOrderUseCase)(nil).DeleteOrder), ctx, id)
}

// GetOrder mocks base method.
func (m *MockOrderUseCase) GetOrder(ctx context.Context, id string) (*usecase.OrderDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", ctx, id)
	ret0, _ := ret[0].(*usecase.OrderDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockOrderUseCaseMockRecorder) GetOrder(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockOrderUseCase)(nil).GetOrder), ctx, id)
}

// GetOrderCreationResources mocks base method.
func (m *MockOrderUseCase) GetOrderCreationResources(ctx context.Context) ([]entity.Customer, []entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderCreationResources", ctx)
	ret0, _ := ret[0].([]entity.Customer)
	ret1, _ := ret[1].([]entity.CatalogItem)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetOrderCreationResources indicates an expected call of GetOrderCreationResources.
func (mr *MockOrderUseCaseMockRecorder) GetOrderCreationResources(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderCreationResources", reflect.TypeOf((*MockOrderUseCase)(nil).GetOrderCreationResources), ctx)
}

// ListOrders mocks base method.
func (m *MockOrderUseCase) ListOrders(ctx context.Context) ([]*usecase.OrderDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrders", ctx)
	ret0, _ := ret[0].([]*usecase.OrderDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOrders indicates an expected call of ListOrders.
func (mr *MockOrderUseCaseMockRecorder) ListOrders(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrders", reflect.TypeOf((*MockOrderUseCase)(nil).ListOrders), ctx)
}
