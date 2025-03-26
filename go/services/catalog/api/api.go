package api

import (
	"github.com/tusmasoma/go-microservice-k8s/go/services/catalog/database"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
)

type catalogService struct {
	db *database.Database
	catalog.UnimplementedCatalogServiceServer
}

func NewCatalogService(db *database.Database) catalog.CatalogServiceServer {
	return &catalogService{
		db: db,
	}
}
