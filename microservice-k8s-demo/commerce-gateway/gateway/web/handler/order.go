package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type OrderHandler interface {
	CreateOrderForm(c *gin.Context)
	CreateOrder(c *gin.Context)
}

type orderHandler struct {
	client pb.OrderServiceClient
}

func NewOrderHandler(client pb.OrderServiceClient) OrderHandler {
	return &orderHandler{
		client: client,
	}
}

func (oh *orderHandler) CreateOrderForm(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := oh.client.GetOrderCreationResources(ctx, &pb.GetOrderCreationResourcesRequest{})
	if err != nil {
		log.Error("Failed to get order page data", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	customers := make([]CustomerData, 0, len(resp.GetCustomers()))
	for _, c := range resp.GetCustomers() {
		customers = append(customers, CustomerData{
			ID:      c.GetId(),
			Name:    c.GetName(),
			Email:   c.GetEmail(),
			Street:  c.GetStreet(),
			City:    c.GetCity(),
			Country: c.GetCountry(),
		})
	}

	items := make([]CatalogItemData, 0, len(resp.GetItems()))
	for _, i := range resp.GetItems() {
		items = append(items, CatalogItemData{
			ID:    i.GetId(),
			Name:  i.GetName(),
			Price: i.GetPrice(),
		})
	}

	c.HTML(http.StatusOK, "order/create.html", gin.H{
		"Customers": customers,
		"Items":     items,
	})
}

type CreateOrderRequest struct {
	CustomerID string `form:"customer_id"`
	Count      int    `form:"count"`
	ItemID     string `form:"item_id"`
}

func (oh *orderHandler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("Failed to bind request", log.Ferror(err))
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	// 実装を簡単にする為に、一旦、一回のorder作成では一つの商品の注文しかできないようにする
	// 上記の制約を受けるのはこの関数のみで、実際のorder作成のロジックには影響しない
	orderLines := []*pb.OrderLine{
		{
			Count:  int32(req.Count),
			ItemId: req.ItemID,
		},
	}

	if _, err := oh.client.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId: req.CustomerID,
		OrderLines: orderLines,
	}); err != nil {
		log.Error("Failed to create order", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/order/list")
}
