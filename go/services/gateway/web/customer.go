package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/request"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/response"
	"github.com/tusmasoma/go-microservice-k8s/go/services/gateway/web/entity"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	webpb "github.com/tusmasoma/go-microservice-k8s/proto/web"
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

// listCustomers godoc
// @Summary      List all customers
// @Description  Get a list of all customers
// @Tags         customer
// @Accept       json
// @Produce      json
// @Success      200  {object}  webpb.ListCustomersResponse
// @Failure      500  {object}  web.Error
// @Router       /api/v1/customer [get]
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

// createCustomer godoc
// @Summary      Create a new customer
// @Description  Create a new customer with the provided information
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request body webpb.CreateCustomerRequest true "Customer information"
// @Success      201
// @Failure      400  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/customer [post]
func (w *web) createCustomer(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &webpb.CreateCustomerRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &customer.CreateCustomerRequest{
		Name:        req.GetName(),
		Email:       req.GetEmail(),
		Street:      req.GetStreet(),
		City:        req.GetCity(),
		CountryCode: customer.CountryCode(req.GetCountryCode()),
	}
	if _, err := w.customer.CreateCustomer(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.Created(rw)
}

// getCustomer godoc
// @Summary      Get a customer
// @Description  Get customer information by ID
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        customerID path string true "Customer ID"
// @Success      200  {object}  webpb.GetCustomerResponse
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/customer/{customerID} [get]
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

// updateCustomer godoc
// @Summary      Update a customer
// @Description  Update customer information by ID
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        customerID path string true "Customer ID"
// @Param        request body webpb.UpdateCustomerRequest true "Updated customer information"
// @Success      204
// @Failure      400  {object}  web.Error
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/customer/{customerID} [put]
func (w *web) updateCustomer(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerID := request.Param(r, "customerID")
	req := &webpb.UpdateCustomerRequest{}
	if err := request.Bind(req, r); err != nil {
		response.BadRequest(err, rw, r)
		return
	}
	in := &customer.UpdateCustomerRequest{
		Id:          customerID,
		Name:        req.GetName(),
		Email:       req.GetEmail(),
		Street:      req.GetStreet(),
		City:        req.GetCity(),
		CountryCode: customer.CountryCode(req.GetCountryCode()),
	}
	if _, err := w.customer.UpdateCustomer(ctx, in); err != nil {
		response.Error(err, rw, r)
		return
	}
	response.NoContent(rw)
}

// deleteCustomer godoc
// @Summary      Delete a customer
// @Description  Delete a customer by ID
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        customerID path string true "Customer ID"
// @Success      204
// @Failure      404  {object}  web.Error
// @Failure      500  {object}  web.Error
// @Router       /api/v1/customer/{customerID} [delete]
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
