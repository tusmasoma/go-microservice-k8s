package api

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tusmasoma/go-microservice-k8s/go/services/catalog/database"
	dbmock "github.com/tusmasoma/go-microservice-k8s/go/services/catalog/database/mock"
	"github.com/tusmasoma/go-microservice-k8s/go/services/catalog/entity"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
)

func Test_GetCatalogItem(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCatalogItem(ctrl)
	gomock.InOrder(
		db.EXPECT().Get(ctx, id).Return(nil, errmock),
		db.EXPECT().Get(ctx, id).Return(
			&entity.CatalogItem{
				CatalogItem: catalog.CatalogItem{
					Id:    id,
					Name:  "item",
					Price: 100,
				},
			}, nil),
	)
	service := &catalogService{
		db: &database.Database{
			CatalogItem: db,
		},
	}
	// bad request
	req := &catalog.GetCatalogItemRequest{}
	_, err := service.GetCatalogItem(ctx, req)
	assert.Error(t, err)
	// db error
	req.Id = id
	resp, err := service.GetCatalogItem(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.GetCatalogItem(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_ListCatalogItemsByName(t *testing.T) {
	t.Parallel()
	name := "item"
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCatalogItem(ctrl)
	gomock.InOrder(
		db.EXPECT().ListByName(ctx, name).Return(nil, errmock),
		db.EXPECT().ListByName(ctx, name).Return(
			entity.CatalogItems{
				&entity.CatalogItem{
					CatalogItem: catalog.CatalogItem{
						Id:    uuid.NewString(),
						Name:  "item",
						Price: 100,
					},
				},
			}, nil),
	)
	service := &catalogService{
		db: &database.Database{
			CatalogItem: db,
		},
	}
	// bad request
	req := &catalog.ListCatalogItemsByNameRequest{}
	_, err := service.ListCatalogItemsByName(ctx, req)
	assert.Error(t, err)
	// db error
	req.Name = name
	resp, err := service.ListCatalogItemsByName(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.ListCatalogItemsByName(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_ListCatalogItemsByIDs(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	ids := []string{id}
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCatalogItem(ctrl)
	gomock.InOrder(
		db.EXPECT().ListByIDs(ctx, ids).Return(nil, errmock),
		db.EXPECT().ListByIDs(ctx, ids).Return(
			entity.CatalogItems{
				&entity.CatalogItem{
					CatalogItem: catalog.CatalogItem{
						Id:    id,
						Name:  "item",
						Price: 100,
					},
				},
			}, nil),
	)
	service := &catalogService{
		db: &database.Database{
			CatalogItem: db,
		},
	}
	// bad request
	req := &catalog.ListCatalogItemsByIDsRequest{}
	_, err := service.ListCatalogItemsByIDs(ctx, req)
	assert.Error(t, err)
	// db error
	req.Ids = ids
	resp, err := service.ListCatalogItemsByIDs(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.ListCatalogItemsByIDs(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_ListCatalogItems(t *testing.T) {
	t.Parallel()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCatalogItem(ctrl)
	gomock.InOrder(
		db.EXPECT().List(ctx).Return(nil, errmock),
		db.EXPECT().List(ctx).Return(
			entity.CatalogItems{
				&entity.CatalogItem{
					CatalogItem: catalog.CatalogItem{
						Id:    uuid.NewString(),
						Name:  "item",
						Price: 100,
					},
				},
			}, nil),
	)
	service := &catalogService{
		db: &database.Database{
			CatalogItem: db,
		},
	}
	// db error
	req := &catalog.ListCatalogItemsRequest{}
	resp, err := service.ListCatalogItems(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.ListCatalogItems(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_CreateCatalogItem(t *testing.T) {
	t.Parallel()
	item := &entity.CatalogItem{
		CatalogItem: catalog.CatalogItem{
			Id:    uuid.NewString(),
			Name:  "item",
			Price: 100,
		},
	}
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCatalogItem(ctrl)
	gomock.InOrder(
		db.EXPECT().Create(
			gomock.Any(),
			gomock.Any(),
		).Do(func(_ context.Context, item *entity.CatalogItem) {
			if item.GetName() != "item" {
				t.Errorf("unexpected Name: got %v, want %v", item.GetName(), "item")
			}
			if item.GetPrice() != 100 {
				t.Errorf("unexpected Price: got %v, want %v", item.GetPrice(), 100)
			}
		}).Return(errmock),
		db.EXPECT().Create(
			gomock.Any(),
			gomock.Any(),
		).Do(func(_ context.Context, item *entity.CatalogItem) {
			if item.GetName() != "item" {
				t.Errorf("unexpected Name: got %v, want %v", item.GetName(), "item")
			}
			if item.GetPrice() != 100 {
				t.Errorf("unexpected Price: got %v, want %v", item.GetPrice(), 100)
			}
		}).Return(nil),
	)
	service := &catalogService{
		db: &database.Database{
			CatalogItem: db,
		},
	}
	// bad request
	req := &catalog.CreateCatalogItemRequest{}
	_, err := service.CreateCatalogItem(ctx, req)
	assert.Error(t, err)
	// db error
	req.Name = item.GetName()
	req.Price = item.GetPrice()
	resp, err := service.CreateCatalogItem(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.CreateCatalogItem(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_UpdateCatalogItem(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	item := &entity.CatalogItem{
		CatalogItem: catalog.CatalogItem{
			Id:    id,
			Name:  "item",
			Price: 100,
		},
	}
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCatalogItem(ctrl)
	gomock.InOrder(
		db.EXPECT().Get(ctx, id).Return(nil, errmock),
		db.EXPECT().Get(ctx, id).Return(item, nil),
		db.EXPECT().Update(
			gomock.Any(),
			gomock.Any(),
		).Do(func(_ context.Context, item *entity.CatalogItem) {
			if item.GetName() != "item" {
				t.Errorf("unexpected Name: got %v, want %v", item.GetName(), "item")
			}
			if item.GetPrice() != 100 {
				t.Errorf("unexpected Price: got %v, want %v", item.GetPrice(), 100)
			}
		}).Return(errmock),
		db.EXPECT().Get(ctx, id).Return(item, nil),
		db.EXPECT().Update(
			gomock.Any(),
			gomock.Any(),
		).Do(func(_ context.Context, item *entity.CatalogItem) {
			if item.GetName() != "item" {
				t.Errorf("unexpected Name: got %v, want %v", item.GetName(), "item")
			}
			if item.GetPrice() != 100 {
				t.Errorf("unexpected Price: got %v, want %v", item.GetPrice(), 100)
			}
		}).Return(nil),
	)
	service := &catalogService{
		db: &database.Database{
			CatalogItem: db,
		},
	}
	// bad request
	req := &catalog.UpdateCatalogItemRequest{}
	_, err := service.UpdateCatalogItem(ctx, req)
	assert.Error(t, err)
	// db error (Get)
	req.Id = item.GetId()
	req.Name = item.GetName()
	req.Price = item.GetPrice()
	resp, err := service.UpdateCatalogItem(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// db error (Update)
	resp, err = service.UpdateCatalogItem(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.UpdateCatalogItem(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_DeleteCatalogItem(t *testing.T) {
	t.Parallel()
	id := uuid.NewString()
	errmock := errors.New("errors")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := dbmock.NewMockCatalogItem(ctrl)
	gomock.InOrder(
		db.EXPECT().Delete(ctx, id).Return(errmock),
		db.EXPECT().Delete(ctx, id).Return(nil),
	)
	service := &catalogService{
		db: &database.Database{
			CatalogItem: db,
		},
	}
	// bad request
	req := &catalog.DeleteCatalogItemRequest{}
	_, err := service.DeleteCatalogItem(ctx, req)
	assert.Error(t, err)
	// db error
	req.Id = id
	resp, err := service.DeleteCatalogItem(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	// success
	resp, err = service.DeleteCatalogItem(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
