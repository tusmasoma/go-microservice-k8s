package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/request"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/response"
	"github.com/tusmasoma/go-microservice-k8s/go/services/gateway/web/entity"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	webpb "github.com/tusmasoma/go-microservice-k8s/proto/web"
)

func (w *web) routeOrder(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		// 認証が必要な場合はこのコメントアウトを外す
		// r.Use(w.requireAuth)
		r.Get("/", w.listOrders)
		r.Post("/", w.createOrder)
		r.Get("/{orderID}", w.getOrder)
		r.Delete("/{orderID}", w.deleteOrder)
	})
}

// listOrders godoc
// @Summary      List all orders
// @Description  Get a list of all orders with customer and catalog item details
// @Tags         order
// @Accept       json
// @Produce      json
// @Success      200  {object}  webpb.ListOrdersResponse
// @Failure      500  {object}  web.Error
// @Router       /api/v1/order [get]
func (w *web) listOrders(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderOuts, err := w.order.ListOrders(ctx, &order.ListOrdersRequest{})
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	var orders entity.Orders
	for _, order := range orderOuts.GetOrders() {
		customerIn := &customer.GetCustomerRequest{
			Id: order.GetCustomerId(),
		}
		customerOut, err := w.customer.GetCustomer(ctx, customerIn)
		if err != nil {
			response.Error(err, rw, r)
			return
		}
		var itemIDs []string
		for _, orderLine := range order.GetOrderLines() {
			itemIDs = append(itemIDs, orderLine.GetCatalogItemId())
		}
		catalogIn := &catalog.ListCatalogItemsByIDsRequest{
			Ids: itemIDs,
		}
		catalogOut, err := w.catalog.ListCatalogItemsByIDs(ctx, catalogIn)
		if err != nil {
			response.Error(err, rw, r)
			return
		}
		orders = append(orders, entity.NewOrder(order, customerOut.GetCustomer(), catalogOut.GetItems()))
	}
	body := &webpb.ListOrdersResponse{
		Orders: orders.Proto(),
	}
	response.OK(body, rw, r)
}

// createOrder godoc
// @Summary      Create a new order
// @Description  Create a new order with customer ID and order lines
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request body webpb.CreateOrderRequest true "Order data"
// @Success      201
// @Failure      400  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/order [post]
func (w *web) createOrder(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &webpb.CreateOrderRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	orderLines := make([]*order.OrderLine, len(req.GetOrderLines()))
	for i, orderLine := range req.GetOrderLines() {
		orderLines[i] = &order.OrderLine{
			CatalogItemId: orderLine.GetCatalogItem().GetId(),
			Quantity:      orderLine.GetQuantity(),
		}
	}
	in := &order.CreateOrderRequest{
		CustomerId: req.GetCustomerId(),
		OrderLines: orderLines,
	}
	if _, err := w.order.CreateOrder(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.Created(rw)
}

// getOrder godoc
// @Summary      Get an order
// @Description  Get an order by ID with customer and catalog item details
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        orderID path string true "Order ID"
// @Success      200  {object}  webpb.GetOrderResponse
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/order/{orderID} [get]
func (w *web) getOrder(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := request.Param(r, "orderID")
	orderIn := &order.GetOrderRequest{
		Id: orderID,
	}
	orderOut, err := w.order.GetOrder(ctx, orderIn)
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	customerIn := &customer.GetCustomerRequest{
		Id: orderOut.GetOrder().GetCustomerId(),
	}
	customerOut, err := w.customer.GetCustomer(ctx, customerIn)
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	var itemIDs []string
	for _, orderLine := range orderOut.GetOrder().GetOrderLines() {
		itemIDs = append(itemIDs, orderLine.GetCatalogItemId())
	}
	catalogIn := &catalog.ListCatalogItemsByIDsRequest{
		Ids: itemIDs,
	}
	catalogOut, err := w.catalog.ListCatalogItemsByIDs(ctx, catalogIn)
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	body := &webpb.GetOrderResponse{
		Order: entity.NewOrder(orderOut.GetOrder(), customerOut.GetCustomer(), catalogOut.GetItems()).Proto(),
	}
	response.OK(body, rw, r)
}

// deleteOrder godoc
// @Summary      Delete an order
// @Description  Delete an order by ID
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        orderID path string true "Order ID"
// @Success      204
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/order/{orderID} [delete]
func (w *web) deleteOrder(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := request.Param(r, "orderID")
	in := &order.DeleteOrderRequest{
		Id: orderID,
	}
	if _, err := w.order.DeleteOrder(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.NoContent(rw)
}
