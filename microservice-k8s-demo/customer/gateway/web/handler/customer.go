package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/microservice-k8s-demo/customer/usecase"
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
	cuc usecase.CustomerUsecase
}

func NewCustomerHandler(cuc usecase.CustomerUsecase) CustomerHandler {
	return &customerHandler{
		cuc: cuc,
	}
}

func (ch *customerHandler) ListCustomers(c *gin.Context) {
	ctx := c.Request.Context()

	customers, err := ch.cuc.ListCustomers(ctx)
	if err != nil {
		log.Error("Failed to list customers", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"Customers": customers,
	})
}

func (ch *customerHandler) CreateCustomerForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
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

	params := ch.convertCreateCustomerReqeuestToParams(&req)
	if err := ch.cuc.CreateCustomer(ctx, params); err != nil {
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

func (ch *customerHandler) convertCreateCustomerReqeuestToParams(req *CreateCustomerRequest) *usecase.CreateCustomerParams {
	return &usecase.CreateCustomerParams{
		Name:    req.Name,
		Email:   req.Email,
		Street:  req.Street,
		City:    req.City,
		Country: req.Country,
	}
}

func (ch *customerHandler) UpdateCustomerForm(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	customer, err := ch.cuc.GetCustomer(ctx, id)
	if err != nil {
		log.Error("Failed to get customer", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.HTML(http.StatusOK, "update.html", gin.H{
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

	params := ch.convertUpdateCustomerReqeuestToParams(&req)
	if err := ch.cuc.UpdateCustomer(ctx, params); err != nil {
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

func (ch *customerHandler) convertUpdateCustomerReqeuestToParams(req *UpdateCustomerRequest) *usecase.UpdateCustomerParams {
	return &usecase.UpdateCustomerParams{
		ID:      req.ID,
		Name:    req.Name,
		Email:   req.Email,
		Street:  req.Street,
		City:    req.City,
		Country: req.Country,
	}
}

func (ch *customerHandler) DeleteCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")
	if id == "" {
		log.Warn("ID is required")
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	if err := ch.cuc.DeleteCustomer(ctx, id); err != nil {
		log.Error("Failed to delete customer", log.Ferror(err))
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusFound, "/customer/list")
}
