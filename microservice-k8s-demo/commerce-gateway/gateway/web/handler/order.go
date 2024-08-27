package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type OrderHandler interface {
	CreateOrderForm(c *gin.Context)
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

	resp, err := oh.client.GetOrderPageData(ctx, &pb.GetOrderPageDataRequest{})
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
