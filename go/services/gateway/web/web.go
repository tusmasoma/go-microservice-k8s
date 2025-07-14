package web

import (
	"github.com/go-chi/chi"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
)

type Web interface {
	Register(mux *chi.Mux)
}

type web struct {
	customer customer.CustomerServiceClient
	catalog  catalog.CatalogServiceClient
	order    order.OrderServiceClient
}

func NewWeb(params *HandlerParams) Web {
	return &web{
		customer: params.CustomerClient,
		catalog:  params.CatalogClient,
		order:    params.OrderClient,
	}
}

// @title			Web API
// @version		1.0
// @description	レポジトリのAPIドキュメントです。
// @servers.url	https://api.example.com/web
// @in header
// @name Authorization
func (w *web) Register(mux *chi.Mux) {
	router := func(r chi.Router) {
		r.Route("/customer", w.routeCustomer)
		r.Route("/catalog", w.routeCatalog)
		r.Route("/order", w.routeOrder)
	}
	mux.Route("/web", router)
}

type HandlerParams struct {
	CustomerClient customer.CustomerServiceClient
	CatalogClient  catalog.CatalogServiceClient
	OrderClient    order.OrderServiceClient
}
