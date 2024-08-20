package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/usecase"
)

type CatalogItemHandler interface {
	ListCatalogItems(c *gin.Context)
	// CreateCatalogItem(w http.ResponseWriter, r *http.Request)
}

type catalogItemHandler struct {
	cuc usecase.CatalogItemUseCase
}

func NewCatalogItemHandler(cuc usecase.CatalogItemUseCase) CatalogItemHandler {
	return &catalogItemHandler{
		cuc: cuc,
	}
}

func (ch *catalogItemHandler) ListCatalogItems(c *gin.Context) {
	ctx := c.Request.Context()

	items, err := ch.cuc.ListCatalogItems(ctx)
	if err != nil {
		log.Error("Failed to list catalog items", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}
	data := gin.H{
		"Items": items,
	}

	c.HTML(http.StatusOK, "list.html", data)
}
