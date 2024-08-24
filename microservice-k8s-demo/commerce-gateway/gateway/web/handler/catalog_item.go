package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/proto"
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
	client pb.CatalogServiceClient
}

func NewCatalogItemHandler(client pb.CatalogServiceClient) CatalogItemHandler {
	return &catalogItemHandler{
		client: client,
	}
}

type CatalogItemData struct {
	ID    string
	Name  string
	Price float64
}

func (ch *catalogItemHandler) GetCatalogItemByNameForm(c *gin.Context) {
	c.HTML(http.StatusOK, "catalog/search.html", nil)
}

func (ch *catalogItemHandler) GetCatalogItemByName(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.PostForm("name")
	if name == "" {
		log.Warn("Name is required")
		c.String(http.StatusBadRequest, "Name is required")
		return
	}

	resp, err := ch.client.ListCatalogItemsByName(ctx, &pb.ListCatalogItemsByNameRequest{
		Name: name,
	})
	if err != nil {
		log.Error("Failed to list catalog items by name", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	var items []CatalogItemData
	for _, item := range resp.GetItems() {
		items = append(items, CatalogItemData{
			ID:    item.GetId(),
			Name:  item.GetName(),
			Price: item.GetPrice(),
		})
	}

	c.HTML(http.StatusOK, "catalog/list.html", gin.H{
		"Items": items,
	})
}

func (ch *catalogItemHandler) ListCatalogItems(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := ch.client.ListCatalogItems(ctx, &pb.ListCatalogItemsRequest{})
	if err != nil {
		log.Error("Failed to list catalog items", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	var items []CatalogItemData
	for _, item := range resp.GetItems() {
		items = append(items, CatalogItemData{
			ID:    item.GetId(),
			Name:  item.GetName(),
			Price: item.GetPrice(),
		})
	}

	c.HTML(http.StatusOK, "catalog/list.html", gin.H{
		"Items": items,
	})
}

func (ch *catalogItemHandler) CreateCatalogItemForm(c *gin.Context) {
	c.HTML(http.StatusOK, "catalog/create.html", nil)
}

type CreateCatalogItemRequest struct {
	Name  string  `form:"name"`
	Price float64 `form:"price"`
}

func (ch *catalogItemHandler) CreateCatalogItem(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateCatalogItemRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("Failed to bind request", log.Ferror(err))
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}
	if !ch.isValidCreateCatalogItemRequest(&req) {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if _, err := ch.client.CreateCatalogItem(ctx, &pb.CreateCatalogItemRequest{
		Name:  req.Name,
		Price: req.Price,
	}); err != nil {
		log.Error("Failed to create catalog item", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/catalog/list")
}

func (ch *catalogItemHandler) isValidCreateCatalogItemRequest(req *CreateCatalogItemRequest) bool {
	if req.Name == "" ||
		req.Price <= 0 {
		log.Warn("Invalid request body: %v", req)
		return false
	}
	return true
}

func (ch *catalogItemHandler) UpdateCatalogItemForm(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "ID is required")
		return
	}

	resp, err := ch.client.GetCatalogItem(ctx, &pb.GetCatalogItemRequest{
		Id: id,
	})
	if err != nil {
		log.Error("Failed to get catalog item", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	item := CatalogItemData{
		ID:    resp.GetItem().GetId(),
		Name:  resp.GetItem().GetName(),
		Price: resp.GetItem().GetPrice(),
	}

	c.HTML(http.StatusOK, "catalog/update.html", gin.H{
		"Item": item,
	})
}

type UpdateCatalogItemRequest struct {
	ID    string  `form:"id"`
	Name  string  `form:"name"`
	Price float64 `form:"price"`
}

func (ch *catalogItemHandler) UpdateCatalogItem(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateCatalogItemRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("Failed to bind request", log.Ferror(err))
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}
	if !ch.isValidUpdateCatalogItemRequest(&req) {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if _, err := ch.client.UpdateCatalogItem(ctx, &pb.UpdateCatalogItemRequest{
		Id:    req.ID,
		Name:  req.Name,
		Price: req.Price,
	}); err != nil {
		log.Error("Failed to update catalog item", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/catalog/list")
}

func (ch *catalogItemHandler) isValidUpdateCatalogItemRequest(req *UpdateCatalogItemRequest) bool {
	if req.ID == "" ||
		req.Name == "" ||
		req.Price <= 0 {
		log.Warn("Invalid request body: %v", req)
		return false
	}
	return true
}

func (ch *catalogItemHandler) DeleteCatalogItem(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "ID is required")
		return
	}

	if _, err := ch.client.DeleteCatalogItem(ctx, &pb.DeleteCatalogItemRequest{
		Id: id,
	}); err != nil {
		log.Error("Failed to delete catalog item", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/catalog/list")
}
