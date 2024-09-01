package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type OrderHandler interface {
	ListOrders(c *gin.Context)
	CreateOrderForm(c *gin.Context)
	CreateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
}

type orderHandler struct {
	client pb.OrderServiceClient
}

func NewOrderHandler(client pb.OrderServiceClient) OrderHandler {
	return &orderHandler{
		client: client,
	}
}

func (oh *orderHandler) ListOrders(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := oh.client.ListOrders(ctx, &pb.ListOrdersRequest{})
	if err != nil {
		log.Error("Failed to get order list", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.HTML(http.StatusOK, "order/list.html", gin.H{
		"Orders": resp.GetOrders(),
	})
}

func (oh *orderHandler) CreateOrderForm(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := oh.client.GetOrderCreationResources(ctx, &pb.GetOrderCreationResourcesRequest{})
	if err != nil {
		log.Error("Failed to get order page data", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.HTML(http.StatusOK, "order/create.html", gin.H{
		"Customers": resp.GetCustomers(),
		"Items":     resp.GetItems(),
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
			Count: int32(req.Count),
			Item: &pb.CatalogItem{
				Id: req.ItemID,
			},
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

func (oh *orderHandler) DeleteOrder(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if _, err := oh.client.DeleteOrder(ctx, &pb.DeleteOrderRequest{
		OrderId: id,
	}); err != nil {
		log.Error("Failed to delete order", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/order/list")
}
