package gateway

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tusmasoma/go-microservice-k8s/services/catalog/database"
	"github.com/tusmasoma/go-microservice-k8s/services/catalog/entity"

	pb "github.com/tusmasoma/go-microservice-k8s/proto/catalog"
)

type CatalogItemHandler interface {
	GetCatalogItem(ctx context.Context, req *pb.GetCatalogItemRequest) (*pb.GetCatalogItemResponse, error)
	ListCatalogItemsByName(ctx context.Context, req *pb.ListCatalogItemsByNameRequest) (*pb.ListCatalogItemsByNameResponse, error)
	ListCatalogItemsByIDs(ctx context.Context, req *pb.ListCatalogItemsByIDsRequest) (*pb.ListCatalogItemsByIDsResponse, error)
	ListCatalogItems(ctx context.Context, req *pb.ListCatalogItemsRequest) (*pb.ListCatalogItemsResponse, error)
	CreateCatalogItem(ctx context.Context, req *pb.CreateCatalogItemRequest) (*pb.CreateCatalogItemResponse, error)
	UpdateCatalogItem(ctx context.Context, req *pb.UpdateCatalogItemRequest) (*pb.UpdateCatalogItemResponse, error)
	DeleteCatalogItem(ctx context.Context, req *pb.DeleteCatalogItemRequest) (*pb.DeleteCatalogItemResponse, error)
}

type catalogItemHandler struct {
	db *database.Database
	pb.UnimplementedCatalogServiceServer
}

func NewCatalogItemHandler(db *database.Database) pb.CatalogServiceServer {
	return &catalogItemHandler{
		db: db,
	}
}

func (ch *catalogItemHandler) GetCatalogItem(ctx context.Context, req *pb.GetCatalogItemRequest) (*pb.GetCatalogItemResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	item, err := ch.db.CatalogItem.Get(ctx, id)
	if err != nil {
		log.Error("Failed to get catalog item", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to get catalog item")
	}
	return &pb.GetCatalogItemResponse{
		Item: item.Proto(),
	}, nil
}

func (ch *catalogItemHandler) ListCatalogItemsByName(ctx context.Context, req *pb.ListCatalogItemsByNameRequest) (*pb.ListCatalogItemsByNameResponse, error) {
	name := req.GetName()
	if name == "" {
		log.Warn("Name is required")
		return nil, status.Errorf(codes.InvalidArgument, "Name is required")
	}
	items, err := ch.db.CatalogItem.ListByName(ctx, name)
	if err != nil {
		log.Error("Failed to list catalog items by name", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by name")
	}
	return &pb.ListCatalogItemsByNameResponse{
		Items: items.Proto(),
	}, nil
}

func (ch *catalogItemHandler) ListCatalogItemsByIDs(ctx context.Context, req *pb.ListCatalogItemsByIDsRequest) (*pb.ListCatalogItemsByIDsResponse, error) {
	ids := req.GetIds()
	if len(ids) == 0 {
		log.Warn("IDs are required")
		return nil, status.Errorf(codes.InvalidArgument, "IDs are required")
	}
	items, err := ch.db.CatalogItem.ListByIDs(ctx, ids)
	if err != nil {
		log.Error("Failed to list catalog items by IDs", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by IDs")
	}
	return &pb.ListCatalogItemsByIDsResponse{
		Items: items.Proto(),
	}, nil
}

func (ch *catalogItemHandler) ListCatalogItems(ctx context.Context, _ *pb.ListCatalogItemsRequest) (*pb.ListCatalogItemsResponse, error) {
	items, err := ch.db.CatalogItem.List(ctx)
	if err != nil {
		log.Error("Failed to list catalog items by name", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by name")
	}
	return &pb.ListCatalogItemsResponse{
		Items: items.Proto(),
	}, nil
}

func (ch *catalogItemHandler) CreateCatalogItem(ctx context.Context, req *pb.CreateCatalogItemRequest) (*pb.CreateCatalogItemResponse, error) {
	if !ch.isValidCreateCatalogItemRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	item, err := entity.CreateCatalogItem(req.GetName(), req.GetPrice())
	if err != nil {
		return nil, err
	}
	if err := ch.db.CatalogItem.Create(ctx, item); err != nil {
		log.Error("Failed to create catalog item", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to create catalog item")
	}
	return &pb.CreateCatalogItemResponse{}, nil
}

func (ch *catalogItemHandler) isValidCreateCatalogItemRequest(req *pb.CreateCatalogItemRequest) bool {
	if req.GetName() == "" ||
		req.GetPrice() <= 0 {
		log.Warn(
			"Invalid request",
			log.Fstring("name", req.GetName()),
			log.Ffloat64("price", req.GetPrice()),
		)
		return false
	}
	return true
}

func (ch *catalogItemHandler) UpdateCatalogItem(ctx context.Context, req *pb.UpdateCatalogItemRequest) (*pb.UpdateCatalogItemResponse, error) {
	if !ch.isValidUpdateCatalogItemRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	item, err := ch.db.CatalogItem.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	item.Name = req.GetName()
	item.Price = req.GetPrice()
	if err := ch.db.CatalogItem.Update(ctx, item); err != nil {
		log.Error("Failed to update catalog item", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to update catalog item")
	}
	return &pb.UpdateCatalogItemResponse{}, nil
}

func (ch *catalogItemHandler) isValidUpdateCatalogItemRequest(req *pb.UpdateCatalogItemRequest) bool {
	if req.GetId() == "" ||
		req.GetName() == "" ||
		req.GetPrice() <= 0 {
		log.Warn(
			"Invalid request",
			log.Fstring("id", req.GetId()),
			log.Fstring("name", req.GetName()),
			log.Ffloat64("price", req.GetPrice()),
		)
		return false
	}
	return true
}

func (ch *catalogItemHandler) DeleteCatalogItem(ctx context.Context, req *pb.DeleteCatalogItemRequest) (*pb.DeleteCatalogItemResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}
	if err := ch.db.CatalogItem.Delete(ctx, id); err != nil {
		log.Error("Failed to delete catalog item", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete catalog item")
	}
	return &pb.DeleteCatalogItemResponse{}, nil
}
