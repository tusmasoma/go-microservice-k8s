package catalogservice

import (
	"context"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/proto"
)

type catalogItemRepository struct {
	client pb.CatalogServiceClient
}

func NewCatalogItemRepository(client pb.CatalogServiceClient) repository.CatalogItemRepository {
	return &catalogItemRepository{
		client,
	}
}

func (r *catalogItemRepository) Get(ctx context.Context, id string) (*entity.CatalogItem, error) {
	resp, err := r.client.GetCatalogItem(ctx, &pb.GetCatalogItemRequest{Id: id})
	if err != nil {
		return nil, err
	}

	item := entity.CatalogItem{
		ID:    resp.GetItem().GetId(),
		Name:  resp.GetItem().GetName(),
		Price: resp.GetItem().GetPrice(),
	}

	return &item, err
}

func (r *catalogItemRepository) List(ctx context.Context) ([]entity.CatalogItem, error) {
	resp, err := r.client.ListCatalogItems(ctx, &pb.ListCatalogItemsRequest{})
	if err != nil {
		return nil, err
	}

	var items []entity.CatalogItem
	for _, i := range resp.GetItems() {
		items = append(items, entity.CatalogItem{
			ID:    i.GetId(),
			Name:  i.GetName(),
			Price: i.GetPrice(),
		})
	}

	return items, nil
}

func (r *catalogItemRepository) ListByName(ctx context.Context, name string) ([]entity.CatalogItem, error) {
	resp, err := r.client.ListCatalogItemsByName(ctx, &pb.ListCatalogItemsByNameRequest{Name: name})
	if err != nil {
		return nil, err
	}

	var items []entity.CatalogItem
	for _, i := range resp.GetItems() {
		items = append(items, entity.CatalogItem{
			ID:    i.GetId(),
			Name:  i.GetName(),
			Price: i.GetPrice(),
		})
	}

	return items, nil
}

func (r *catalogItemRepository) Create(ctx context.Context, item entity.CatalogItem) error {
	if _, err := r.client.CreateCatalogItem(ctx, &pb.CreateCatalogItemRequest{
		Name:  item.Name,
		Price: item.Price,
	}); err != nil {
		return err
	}
	return nil
}

func (r *catalogItemRepository) Update(ctx context.Context, item entity.CatalogItem) error {
	if _, err := r.client.UpdateCatalogItem(ctx, &pb.UpdateCatalogItemRequest{
		Id:    item.ID,
		Name:  item.Name,
		Price: item.Price,
	}); err != nil {
		return err
	}
	return nil
}

func (r *catalogItemRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.client.DeleteCatalogItem(ctx, &pb.DeleteCatalogItemRequest{Id: id}); err != nil {
		return err
	}
	return nil
}
