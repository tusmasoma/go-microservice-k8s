// Code generated by MockGen. DO NOT EDIT.
// Source: catalog_item.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/tusmasoma/go-microservice-k8s/services/catalog/entity"
)

// MockCatalogItemUseCase is a mock of CatalogItemUseCase interface.
type MockCatalogItemUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCatalogItemUseCaseMockRecorder
}

// MockCatalogItemUseCaseMockRecorder is the mock recorder for MockCatalogItemUseCase.
type MockCatalogItemUseCaseMockRecorder struct {
	mock *MockCatalogItemUseCase
}

// NewMockCatalogItemUseCase creates a new mock instance.
func NewMockCatalogItemUseCase(ctrl *gomock.Controller) *MockCatalogItemUseCase {
	mock := &MockCatalogItemUseCase{ctrl: ctrl}
	mock.recorder = &MockCatalogItemUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCatalogItemUseCase) EXPECT() *MockCatalogItemUseCaseMockRecorder {
	return m.recorder
}

// CreateCatalogItem mocks base method.
func (m *MockCatalogItemUseCase) CreateCatalogItem(ctx context.Context, name string, price float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCatalogItem", ctx, name, price)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCatalogItem indicates an expected call of CreateCatalogItem.
func (mr *MockCatalogItemUseCaseMockRecorder) CreateCatalogItem(ctx, name, price interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCatalogItem", reflect.TypeOf((*MockCatalogItemUseCase)(nil).CreateCatalogItem), ctx, name, price)
}

// DeleteCatalogItem mocks base method.
func (m *MockCatalogItemUseCase) DeleteCatalogItem(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCatalogItem", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCatalogItem indicates an expected call of DeleteCatalogItem.
func (mr *MockCatalogItemUseCaseMockRecorder) DeleteCatalogItem(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCatalogItem", reflect.TypeOf((*MockCatalogItemUseCase)(nil).DeleteCatalogItem), ctx, id)
}

// GetCatalogItem mocks base method.
func (m *MockCatalogItemUseCase) GetCatalogItem(ctx context.Context, id string) (*entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogItem", ctx, id)
	ret0, _ := ret[0].(*entity.CatalogItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogItem indicates an expected call of GetCatalogItem.
func (mr *MockCatalogItemUseCaseMockRecorder) GetCatalogItem(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogItem", reflect.TypeOf((*MockCatalogItemUseCase)(nil).GetCatalogItem), ctx, id)
}

// ListCatalogItems mocks base method.
func (m *MockCatalogItemUseCase) ListCatalogItems(ctx context.Context) ([]entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCatalogItems", ctx)
	ret0, _ := ret[0].([]entity.CatalogItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCatalogItems indicates an expected call of ListCatalogItems.
func (mr *MockCatalogItemUseCaseMockRecorder) ListCatalogItems(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCatalogItems", reflect.TypeOf((*MockCatalogItemUseCase)(nil).ListCatalogItems), ctx)
}

// ListCatalogItemsByIDs mocks base method.
func (m *MockCatalogItemUseCase) ListCatalogItemsByIDs(ctx context.Context, ids []string) ([]entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCatalogItemsByIDs", ctx, ids)
	ret0, _ := ret[0].([]entity.CatalogItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCatalogItemsByIDs indicates an expected call of ListCatalogItemsByIDs.
func (mr *MockCatalogItemUseCaseMockRecorder) ListCatalogItemsByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCatalogItemsByIDs", reflect.TypeOf((*MockCatalogItemUseCase)(nil).ListCatalogItemsByIDs), ctx, ids)
}

// ListCatalogItemsByName mocks base method.
func (m *MockCatalogItemUseCase) ListCatalogItemsByName(ctx context.Context, name string) ([]entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCatalogItemsByName", ctx, name)
	ret0, _ := ret[0].([]entity.CatalogItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCatalogItemsByName indicates an expected call of ListCatalogItemsByName.
func (mr *MockCatalogItemUseCaseMockRecorder) ListCatalogItemsByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCatalogItemsByName", reflect.TypeOf((*MockCatalogItemUseCase)(nil).ListCatalogItemsByName), ctx, name)
}

// UpdateCatalogItem mocks base method.
func (m *MockCatalogItemUseCase) UpdateCatalogItem(ctx context.Context, id, name string, price float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCatalogItem", ctx, id, name, price)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCatalogItem indicates an expected call of UpdateCatalogItem.
func (mr *MockCatalogItemUseCaseMockRecorder) UpdateCatalogItem(ctx, id, name, price interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCatalogItem", reflect.TypeOf((*MockCatalogItemUseCase)(nil).UpdateCatalogItem), ctx, id, name, price)
}
