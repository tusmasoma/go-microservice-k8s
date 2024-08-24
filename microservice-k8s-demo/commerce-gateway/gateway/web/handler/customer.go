package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/commerce-gateway/entity"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/proto"
)

type CustomerHandler interface {
	ListCustomers(c *gin.Context)
	CreateCustomerForm(c *gin.Context)
	CreateCustomer(c *gin.Context)
	UpdateCustomerForm(c *gin.Context)
	UpdateCustomer(c *gin.Context)
	DeleteCustomer(c *gin.Context)
}

type customerHandler struct {
	client pb.CustomerServiceClient
}

func NewCustomerHandler(client pb.CustomerServiceClient) CustomerHandler {
	return &customerHandler{
		client: client,
	}
}

func (ch *customerHandler) ListCustomers(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := ch.client.ListCustomers(ctx, &pb.ListCustomersRequest{})
	if err != nil {
		log.Error("Failed to list customers", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	var customers []entity.Customer
	for _, c := range resp.GetCustomers() {
		customers = append(customers, entity.Customer{
			ID:      c.GetId(),
			Name:    c.GetName(),
			Email:   c.GetEmail(),
			Street:  c.GetStreet(),
			City:    c.GetCity(),
			Country: c.GetCountry(),
		})
	}

	c.HTML(http.StatusOK, "customer/list.html", gin.H{
		"Customers": customers,
	})
}

func (ch *customerHandler) CreateCustomerForm(c *gin.Context) {
	c.HTML(http.StatusOK, "customer/create.html", nil)
}

type CreateCustomerRequest struct {
	Name    string `form:"name"`
	Email   string `form:"email"`
	Street  string `form:"street"`
	City    string `form:"city"`
	Country string `form:"country"`
}

func (ch *customerHandler) CreateCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateCustomerRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("Failed to bind request", log.Ferror(err))
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if !ch.isValidCreateCustomerRequest(&req) {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if _, err := ch.client.CreateCustomer(ctx, &pb.CreateCustomerRequest{
		Name:    req.Name,
		Email:   req.Email,
		Street:  req.Street,
		City:    req.City,
		Country: req.Country,
	}); err != nil {
		log.Error("Failed to create customer", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/customer/list")
}

func (ch *customerHandler) isValidCreateCustomerRequest(req *CreateCustomerRequest) bool {
	if req.Name == "" ||
		req.Email == "" ||
		req.Street == "" ||
		req.City == "" ||
		req.Country == "" {
		log.Warn("Invalid request body: %v", req)
		return false
	}
	return true
}

func (ch *customerHandler) UpdateCustomerForm(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := ch.client.GetCustomer(ctx, &pb.GetCustomerRequest{Id: id})
	if err != nil {
		log.Error("Failed to get customer", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	customer := entity.Customer{
		ID:      resp.GetCustomer().GetId(),
		Name:    resp.GetCustomer().GetName(),
		Email:   resp.GetCustomer().GetEmail(),
		Street:  resp.GetCustomer().GetStreet(),
		City:    resp.GetCustomer().GetCity(),
		Country: resp.GetCustomer().GetCountry(),
	}

	c.HTML(http.StatusOK, "customer/update.html", gin.H{
		"Customer": customer,
	})
}

type UpdateCustomerRequest struct {
	ID      string `form:"id"`
	Name    string `form:"name"`
	Email   string `form:"email"`
	Street  string `form:"street"`
	City    string `form:"city"`
	Country string `form:"country"`
}

func (ch *customerHandler) UpdateCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateCustomerRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("Failed to bind request", log.Ferror(err))
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if !ch.isValidUpdateCustomerRequest(&req) {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if _, err := ch.client.UpdateCustomer(ctx, &pb.UpdateCustomerRequest{
		Id:      req.ID,
		Name:    req.Name,
		Email:   req.Email,
		Street:  req.Street,
		City:    req.City,
		Country: req.Country,
	}); err != nil {
		log.Error("Failed to update customer", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/customer/list")
}

func (ch *customerHandler) isValidUpdateCustomerRequest(req *UpdateCustomerRequest) bool {
	if req.ID == "" ||
		req.Name == "" ||
		req.Email == "" ||
		req.Street == "" ||
		req.City == "" ||
		req.Country == "" {
		log.Warn("Invalid request body: %v", req)
		return false
	}
	return true
}

func (ch *customerHandler) DeleteCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if _, err := ch.client.DeleteCustomer(ctx, &pb.DeleteCustomerRequest{Id: id}); err != nil {
		log.Error("Failed to delete customer", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/customer/list")
}
