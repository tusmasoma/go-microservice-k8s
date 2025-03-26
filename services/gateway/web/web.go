package web

import (
	"github.com/go-chi/chi"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
)

type web struct {
	customer customer.CustomerServiceClient
	catalog  catalog.CatalogServiceClient
	order    order.OrderServiceClient
}

func (w *web) Register(mux *chi.Mux) {
	router := func(r chi.Router) {
		r.Route("/customer", w.routeCustomer)
		r.Route("/catalog", w.routeCatalog)
		r.Route("/order", w.routeOrder)
	}
	mux.Route("/web", router)
}
