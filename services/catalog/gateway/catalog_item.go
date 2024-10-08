package gateway

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tusmasoma/go-microservice-k8s/services/catalog/usecase"

	pb "github.com/tusmasoma/go-microservice-k8s/services/catalog/proto"
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
	cuc usecase.CatalogItemUseCase
	pb.UnimplementedCatalogServiceServer
}

func NewCatalogItemHandler(cuc usecase.CatalogItemUseCase) pb.CatalogServiceServer {
	return &catalogItemHandler{
		cuc: cuc,
	}
}

func (ch *catalogItemHandler) GetCatalogItem(ctx context.Context, req *pb.GetCatalogItemRequest) (*pb.GetCatalogItemResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}

	item, err := ch.cuc.GetCatalogItem(ctx, id)
	if err != nil {
		log.Error("Failed to get catalog item", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to get catalog item")
	}

	return &pb.GetCatalogItemResponse{
		Item: &pb.CatalogItem{
			Id:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		},
	}, nil
}

func (ch *catalogItemHandler) ListCatalogItemsByName(ctx context.Context, req *pb.ListCatalogItemsByNameRequest) (*pb.ListCatalogItemsByNameResponse, error) {
	name := req.GetName()
	if name == "" {
		log.Warn("Name is required")
		return nil, status.Errorf(codes.InvalidArgument, "Name is required")
	}

	items, err := ch.cuc.ListCatalogItemsByName(ctx, name)
	if err != nil {
		log.Error("Failed to list catalog items by name", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by name")
	}

	var res []*pb.CatalogItem
	for _, item := range items {
		res = append(res, &pb.CatalogItem{
			Id:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &pb.ListCatalogItemsByNameResponse{
		Items: res,
	}, nil
}

func (ch *catalogItemHandler) ListCatalogItemsByIDs(ctx context.Context, req *pb.ListCatalogItemsByIDsRequest) (*pb.ListCatalogItemsByIDsResponse, error) {
	ids := req.GetIds()
	if len(ids) == 0 {
		log.Warn("IDs are required")
		return nil, status.Errorf(codes.InvalidArgument, "IDs are required")
	}

	items, err := ch.cuc.ListCatalogItemsByIDs(ctx, ids)
	if err != nil {
		log.Error("Failed to list catalog items by IDs", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by IDs")
	}

	var res []*pb.CatalogItem
	for _, item := range items {
		res = append(res, &pb.CatalogItem{
			Id:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &pb.ListCatalogItemsByIDsResponse{
		Items: res,
	}, nil
}

func (ch *catalogItemHandler) ListCatalogItems(ctx context.Context, _ *pb.ListCatalogItemsRequest) (*pb.ListCatalogItemsResponse, error) {
	items, err := ch.cuc.ListCatalogItems(ctx)
	if err != nil {
		log.Error("Failed to list catalog items by name", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list catalog items by name")
	}

	var res []*pb.CatalogItem
	for _, item := range items {
		res = append(res, &pb.CatalogItem{
			Id:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &pb.ListCatalogItemsResponse{
		Items: res,
	}, nil
}

func (ch *catalogItemHandler) CreateCatalogItem(ctx context.Context, req *pb.CreateCatalogItemRequest) (*pb.CreateCatalogItemResponse, error) {
	if !ch.isValidCreateCatalogItemRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	if err := ch.cuc.CreateCatalogItem(
		ctx,
		req.GetName(),
		req.GetPrice(),
	); err != nil {
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

	if err := ch.cuc.UpdateCatalogItem(
		ctx,
		req.GetId(),
		req.GetName(),
		req.GetPrice(),
	); err != nil {
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

	if err := ch.cuc.DeleteCatalogItem(ctx, id); err != nil {
		log.Error("Failed to delete catalog item", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete catalog item")
	}

	return &pb.DeleteCatalogItemResponse{}, nil
}
