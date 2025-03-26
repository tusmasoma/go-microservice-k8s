package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tusmasoma/go-microservice-k8s/pkg/request"
	"github.com/tusmasoma/go-microservice-k8s/pkg/response"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	webpb "github.com/tusmasoma/go-microservice-k8s/proto/web"
	"github.com/tusmasoma/go-microservice-k8s/services/gateway/web/entity"
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
		for _, ol := range order.GetOrderLines() {
			itemIDs = append(itemIDs, ol.GetCatalogItemId())
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

func (w *web) createOrder(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &webpb.CreateOrderRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &order.CreateOrderRequest{
		CustomerId: req.GetCustomerId(),
		OrderLines: nil, // TODO: 後で実装
	}
	if _, err := w.order.CreateOrder(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.Created(rw)
}

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
	for _, ol := range orderOut.GetOrder().GetOrderLines() {
		itemIDs = append(itemIDs, ol.GetCatalogItemId())
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
