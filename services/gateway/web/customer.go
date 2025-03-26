package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tusmasoma/go-microservice-k8s/pkg/request"
	"github.com/tusmasoma/go-microservice-k8s/pkg/response"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	webpb "github.com/tusmasoma/go-microservice-k8s/proto/web"
	"github.com/tusmasoma/go-microservice-k8s/services/gateway/web/entity"
)

func (w *web) routeCustomer(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		// 認証が必要な場合はこのコメントアウトを外す
		// r.Use(w.requireAuth)
		r.Get("/", w.listCustomers)
		r.Post("/", w.createCustomer)
		r.Get("/{customerID}", w.getCustomer)
		r.Put("/{customerID}", w.updateCustomer)
		r.Delete("/{customerID}", w.deleteCustomer)
	})
}

func (w *web) listCustomers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := w.customer.ListCustomers(ctx, &customer.ListCustomersRequest{})
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	body := &webpb.ListCustomersResponse{
		Customers: entity.NewCustomers(out.GetCustomers()).Proto(),
	}
	response.OK(body, rw, r)
}

func (w *web) createCustomer(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &webpb.CreateCustomerRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &customer.CreateCustomerRequest{
		Name:    req.GetName(),
		Email:   req.GetEmail(),
		Street:  req.GetStreet(),
		City:    req.GetCity(),
		Country: req.GetCountry(),
	}
	if _, err := w.customer.CreateCustomer(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.Created(rw)
}

func (w *web) getCustomer(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerID := request.Param(r, "customerID")
	in := &customer.GetCustomerRequest{
		Id: customerID,
	}
	out, err := w.customer.GetCustomer(ctx, in)
	if err != nil {
		response.Error(err, rw, r)
		return
	}
	body := &webpb.GetCustomerResponse{
		Customer: entity.NewCustomer(out.GetCustomer()).Proto(),
	}
	response.OK(body, rw, r)
}

func (w *web) updateCustomer(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerID := request.Param(r, "customerID")
	req := &webpb.UpdateCustomerRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &customer.UpdateCustomerRequest{
		Id:      customerID,
		Name:    req.GetName(),
		Email:   req.GetEmail(),
		Street:  req.GetStreet(),
		City:    req.GetCity(),
		Country: req.GetCountry(),
	}
	if _, err := w.customer.UpdateCustomer(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.NoContent(rw)
}

func (w *web) deleteCustomer(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerID := request.Param(r, "customerID")
	in := &customer.DeleteCustomerRequest{
		Id: customerID,
	}
	if _, err := w.customer.DeleteCustomer(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.NoContent(rw)
}
