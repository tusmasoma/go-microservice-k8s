package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tusmasoma/go-microservice-k8s/go/services/catalog/entity"

	pb "github.com/tusmasoma/go-microservice-k8s/proto/catalog"
)

func (c *catalogService) GetCatalogItem(ctx context.Context, req *pb.GetCatalogItemRequest) (*pb.GetCatalogItemResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	item, err := c.db.CatalogItem.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get catalog item")
	}
	return &pb.GetCatalogItemResponse{
		Item: item.Proto(),
	}, nil
}

func (c *catalogService) ListCatalogItemsByName(ctx context.Context, req *pb.ListCatalogItemsByNameRequest) (*pb.ListCatalogItemsByNameResponse, error) {
	name := req.GetName()
	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name is required")
	}
	items, err := c.db.CatalogItem.ListByName(ctx, name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by name")
	}
	return &pb.ListCatalogItemsByNameResponse{
		Items: items.Proto(),
	}, nil
}

func (c *catalogService) ListCatalogItemsByIDs(ctx context.Context, req *pb.ListCatalogItemsByIDsRequest) (*pb.ListCatalogItemsByIDsResponse, error) {
	ids := req.GetIds()
	if len(ids) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "IDs are required")
	}
	items, err := c.db.CatalogItem.ListByIDs(ctx, ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by IDs")
	}
	return &pb.ListCatalogItemsByIDsResponse{
		Items: items.Proto(),
	}, nil
}

func (c *catalogService) ListCatalogItems(ctx context.Context, _ *pb.ListCatalogItemsRequest) (*pb.ListCatalogItemsResponse, error) {
	items, err := c.db.CatalogItem.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by name")
	}
	return &pb.ListCatalogItemsResponse{
		Items: items.Proto(),
	}, nil
}

func (c *catalogService) CreateCatalogItem(ctx context.Context, req *pb.CreateCatalogItemRequest) (*pb.CreateCatalogItemResponse, error) {
	if !c.isValidCreateCatalogItemRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	item, err := entity.CreateCatalogItem(req.GetName(), req.GetPrice())
	if err != nil {
		return nil, err
	}
	if err := c.db.CatalogItem.Create(ctx, item); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create catalog item")
	}
	return &pb.CreateCatalogItemResponse{}, nil
}

func (c *catalogService) isValidCreateCatalogItemRequest(req *pb.CreateCatalogItemRequest) bool {
	if req.GetName() == "" ||
		req.GetPrice() <= 0 {
		return false
	}
	return true
}

func (c *catalogService) UpdateCatalogItem(ctx context.Context, req *pb.UpdateCatalogItemRequest) (*pb.UpdateCatalogItemResponse, error) {
	if !c.isValidUpdateCatalogItemRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	item, err := c.db.CatalogItem.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	item.Name = req.GetName()
	item.Price = req.GetPrice()
	if err := c.db.CatalogItem.Update(ctx, item); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update catalog item")
	}
	return &pb.UpdateCatalogItemResponse{}, nil
}

func (c *catalogService) isValidUpdateCatalogItemRequest(req *pb.UpdateCatalogItemRequest) bool {
	if req.GetId() == "" ||
		req.GetName() == "" ||
		req.GetPrice() <= 0 {
		return false
	}
	return true
}

func (c *catalogService) DeleteCatalogItem(ctx context.Context, req *pb.DeleteCatalogItemRequest) (*pb.DeleteCatalogItemResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	if err := c.db.CatalogItem.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete catalog item")
	}
	return &pb.DeleteCatalogItemResponse{}, nil
}
