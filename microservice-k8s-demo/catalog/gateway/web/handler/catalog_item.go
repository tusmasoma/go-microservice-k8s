package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/usecase"
)

type CatalogItemHandler interface {
	ListCatalogItems(c *gin.Context)
	CreateCatalogItem(c *gin.Context)
	DeleteCatalogItem(c *gin.Context)
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

func (ch *catalogItemHandler) CreateCatalogItem(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Error("Failed to parse price", log.Ferror(err))
		c.String(http.StatusBadRequest, "Invalid price format")
		return
	}

	if err = ch.cuc.CreateCatalogItem(ctx, name, price); err != nil {
		log.Error("Failed to create catalog item", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/catalog/list")
}

func (ch *catalogItemHandler) DeleteCatalogItem(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "ID is required")
		return
	}

	if err := ch.cuc.DeleteCatalogItem(ctx, id); err != nil {
		log.Error("Failed to delete catalog item", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/catalog/list")
}
