package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/usecase"
)

type CatalogItemHandler interface {
	GetCatalogItemByNameForm(c *gin.Context)
	GetCatalogItemByName(c *gin.Context)
	ListCatalogItems(c *gin.Context)
	CreateCatalogItemForm(c *gin.Context)
	CreateCatalogItem(c *gin.Context)
	UpdateCatalogItemForm(c *gin.Context)
	UpdateCatalogItem(c *gin.Context)
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

func (ch *catalogItemHandler) GetCatalogItemByNameForm(c *gin.Context) {
	c.HTML(http.StatusOK, "search.html", nil)
}

func (ch *catalogItemHandler) GetCatalogItemByName(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.PostForm("name")
	if name == "" {
		log.Warn("Name is required")
		c.String(http.StatusBadRequest, "Name is required")
		return
	}

	items, err := ch.cuc.ListCatalogItemsByNameContaining(ctx, name)
	if err != nil {
		log.Error("Failed to list catalog items by name", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}
	data := gin.H{
		"Items": items,
	}

	c.HTML(http.StatusOK, "list.html", data)
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

func (ch *catalogItemHandler) CreateCatalogItemForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
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

func (ch *catalogItemHandler) UpdateCatalogItemForm(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "ID is required")
		return
	}

	item, err := ch.cuc.GetCatalogItem(ctx, id)
	if err != nil {
		log.Error("Failed to get catalog item", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	data := gin.H{
		"Item": item,
	}

	c.HTML(http.StatusOK, "update.html", data)
}

func (ch *catalogItemHandler) UpdateCatalogItem(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.PostForm("id")
	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Error("Failed to parse price", log.Ferror(err))
		c.String(http.StatusBadRequest, "Invalid price format")
		return
	}

	if err = ch.cuc.UpdateCatalogItem(ctx, id, name, price); err != nil {
		log.Error("Failed to update catalog item", log.Ferror(err))
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
