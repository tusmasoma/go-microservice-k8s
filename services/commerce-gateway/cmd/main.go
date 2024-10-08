package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tusmasoma/go-microservice-k8s/services/catalog/config"
	"github.com/tusmasoma/go-microservice-k8s/services/commerce-gateway/gateway/web/handler"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc"

	catalog_pb "github.com/tusmasoma/go-microservice-k8s/services/catalog/proto"
	cusotmer_pb "github.com/tusmasoma/go-microservice-k8s/services/customer/proto"
	order_pb "github.com/tusmasoma/go-microservice-k8s/services/order/proto"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found", log.Ferror(err))
	}

	var addr string
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()

	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	srv, err := BuildContainer(mainCtx, addr)
	if err != nil {
		log.Critical("Failed to build container", log.Ferror(err))
		return
	}
	log.Info("Server running...")

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed", log.Ferror(err))
			return
		}
	}()

	<-signalCtx.Done()
	log.Info("Server stopping...")

	tctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd // 5 is reasonable
	defer cancelShutdown()

	if err = srv.Shutdown(tctx); err != nil {
		log.Error("Failed to shutdown http server", log.Ferror(err))
	}
	log.Info("Server exited")
}

func BuildContainer(ctx context.Context, addr string) (*http.Server, error) { //nolint:funlen // it's okay
	serverConfig, err := config.NewServerConfig(ctx)
	if err != nil {
		log.Critical("Failed to load server config", log.Ferror(err))
		return nil, err
	}

	catalogConn, err := grpc.Dial("catalog-service:8082", grpc.WithInsecure()) //nolint:staticcheck // ignore deprecation
	if err != nil {
		log.Critical("Failed to connect to catalog service", log.Ferror(err))
		return nil, err
	}

	customerConn, err := grpc.Dial("customer-service:8081", grpc.WithInsecure()) //nolint:staticcheck // ignore deprecation
	if err != nil {
		log.Critical("Failed to connect to customer service", log.Ferror(err))
		return nil, err
	}
	orderConn, err := grpc.Dial("order-service:8083", grpc.WithInsecure()) //nolint:staticcheck // ignore deprecation
	if err != nil {
		log.Critical("Failed to connect to order service", log.Ferror(err))
		return nil, err
	}

	catalogClient := catalog_pb.NewCatalogServiceClient(catalogConn)
	customerClient := cusotmer_pb.NewCustomerServiceClient(customerConn)
	orderClient := order_pb.NewOrderServiceClient(orderConn)

	catalogHandler := handler.NewCatalogItemHandler(catalogClient)
	customerHandler := handler.NewCustomerHandler(customerClient)
	orderHandler := handler.NewOrderHandler(orderClient)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposeHeaders:    []string{"Link", "Authorization"},
		AllowCredentials: true,
		MaxAge:           time.Duration(serverConfig.PreflightCacheDurationSec) * time.Second,
	}))

	r.LoadHTMLFiles("gateway/web/templates/index.html")
	r.LoadHTMLGlob("gateway/web/templates/**/*")

	api := r.Group("/")
	{
		api.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "base/index.html", gin.H{})
		})
	}
	{
		catalog := api.Group("/catalog")
		{
			// List all catalog items
			catalog.GET("/list", catalogHandler.ListCatalogItems)

			// Show the form to create a new catalog item
			catalog.GET("/create", catalogHandler.CreateCatalogItemForm)

			// Process the form submission to create a new catalog item
			catalog.POST("/create", catalogHandler.CreateCatalogItem)

			// Show the form to update a catalog item
			catalog.GET("/update", catalogHandler.UpdateCatalogItemForm)

			// Process the form submission to update a catalog item
			catalog.POST("/update", catalogHandler.UpdateCatalogItem)

			// Delete a catalog item
			catalog.GET("/delete", catalogHandler.DeleteCatalogItem)

			// Show the form to search for catalog items by name
			catalog.GET("/search", catalogHandler.GetCatalogItemByNameForm)

			// Process the form submission to search for catalog items by name
			catalog.POST("/search", catalogHandler.GetCatalogItemByName)
		}
	}
	{
		customer := api.Group("/customer")
		{
			// List all customers
			customer.GET("/list", customerHandler.ListCustomers)

			// Show the form to create a new customer
			customer.GET("/create", customerHandler.CreateCustomerForm)

			// Process the form submission to create a new customer
			customer.POST("/create", customerHandler.CreateCustomer)

			// Show the form to update a customer
			customer.GET("/update", customerHandler.UpdateCustomerForm)

			// Process the form submission to update a customer
			customer.POST("/update", customerHandler.UpdateCustomer)

			// Delete a customer
			customer.GET("/delete", customerHandler.DeleteCustomer)
		}
	}
	{
		order := api.Group("/order")
		{
			// List all orders
			order.GET("/list", orderHandler.ListOrders)

			// Show the form to create a new order
			order.GET("/create", orderHandler.CreateOrderForm)

			// Process the form submission to create a new order
			order.POST("/create", orderHandler.CreateOrder)

			// Delete an order
			order.GET("/delete", orderHandler.DeleteOrder)
		}
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}

	return srv, nil
}
