// Code generated by MockGen. DO NOT EDIT.
// Source: catalog_item.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/tusmasoma/microservice-k8s-demo/catalog/entity"
)

// MockCatalogItemRepository is a mock of CatalogItemRepository interface.
type MockCatalogItemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCatalogItemRepositoryMockRecorder
}

// MockCatalogItemRepositoryMockRecorder is the mock recorder for MockCatalogItemRepository.
type MockCatalogItemRepositoryMockRecorder struct {
	mock *MockCatalogItemRepository
}

// NewMockCatalogItemRepository creates a new mock instance.
func NewMockCatalogItemRepository(ctrl *gomock.Controller) *MockCatalogItemRepository {
	mock := &MockCatalogItemRepository{ctrl: ctrl}
	mock.recorder = &MockCatalogItemRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCatalogItemRepository) EXPECT() *MockCatalogItemRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCatalogItemRepository) Create(ctx context.Context, item entity.CatalogItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCatalogItemRepositoryMockRecorder) Create(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCatalogItemRepository)(nil).Create), ctx, item)
}

// Delete mocks base method.
func (m *MockCatalogItemRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCatalogItemRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCatalogItemRepository)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockCatalogItemRepository) Get(ctx context.Context, id string) (*entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.CatalogItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCatalogItemRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCatalogItemRepository)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockCatalogItemRepository) List(ctx context.Context) ([]entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]entity.CatalogItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCatalogItemRepositoryMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCatalogItemRepository)(nil).List), ctx)
}

// ListByName mocks base method.
func (m *MockCatalogItemRepository) ListByName(ctx context.Context, name string) ([]entity.CatalogItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByName", ctx, name)
	ret0, _ := ret[0].([]entity.CatalogItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByName indicates an expected call of ListByName.
func (mr *MockCatalogItemRepositoryMockRecorder) ListByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByName", reflect.TypeOf((*MockCatalogItemRepository)(nil).ListByName), ctx, name)
}

// Update mocks base method.
func (m *MockCatalogItemRepository) Update(ctx context.Context, item entity.CatalogItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCatalogItemRepositoryMockRecorder) Update(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCatalogItemRepository)(nil).Update), ctx, item)
}
