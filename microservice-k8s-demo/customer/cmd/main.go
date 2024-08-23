package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func main() {
}

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Info("No .env file found", log.Ferror(err))
// 	}

// 	var addr string
// 	flag.StringVar(&addr, "addr", ":8081", "tcp host:port to connect")
// 	flag.Parse()

// 	mainCtx, cancelMain := context.WithCancel(context.Background())
// 	defer cancelMain()

// 	container, err := BuildContainer(mainCtx)
// 	if err != nil {
// 		log.Critical("Failed to build container", log.Ferror(err))
// 		return
// 	}

// 	err = container.Invoke(func(router *gin.Engine, config *config.ServerConfig) {
// 		srv := &http.Server{
// 			Addr:         addr,
// 			Handler:      router,
// 			ReadTimeout:  config.ReadTimeout,
// 			WriteTimeout: config.WriteTimeout,
// 			IdleTimeout:  config.IdleTimeout,
// 		}
// 		log.Info("Server running...")

// 		signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
// 		defer stop()

// 		go func() {
// 			if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 				log.Error("Server failed", log.Ferror(err))
// 				return
// 			}
// 		}()

// 		<-signalCtx.Done()
// 		log.Info("Server stopping...")

// 		tctx, cancelShutdown := context.WithTimeout(context.Background(), config.GracefulShutdownTimeout)
// 		defer cancelShutdown()

// 		if err = srv.Shutdown(tctx); err != nil {
// 			log.Error("Failed to shutdown http server", log.Ferror(err))
// 		}
// 		log.Info("Server exited")
// 	})
// 	if err != nil {
// 		log.Critical("Failed to start server", log.Ferror(err))
// 		return
// 	}
// }

// func BuildContainer(ctx context.Context) (*dig.Container, error) {
// 	container := dig.New()

// 	if err := container.Provide(func() context.Context {
// 		return ctx
// 	}); err != nil {
// 		log.Error("Failed to provide context")
// 		return nil, err
// 	}

// 	providers := []interface{}{
// 		config.NewServerConfig,
// 		config.NewDBConfig,
// 		mysql.NewMySQLDB,
// 		mysql.NewTransactionRepository,
// 		mysql.NewCustomerRepository,
// 		usecase.NewCustomerUsecase,
// 		handler.NewCustomerHandler,
// 		func(
// 			serverConfig *config.ServerConfig,
// 			customerHandler handler.CustomerHandler,
// 		) *gin.Engine {
// 			r := gin.Default()

// 			r.Use(cors.New(cors.Config{
// 				AllowOrigins:     []string{"https://*", "http://*"},
// 				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 				AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
// 				ExposeHeaders:    []string{"Link", "Authorization"},
// 				AllowCredentials: true,
// 				MaxAge:           time.Duration(serverConfig.PreflightCacheDurationSec) * time.Second,
// 			}))

// 			r.LoadHTMLGlob("gateway/web/templates/*.html")

// 			api := r.Group("/")
// 			{
// 				customer := api.Group("/customer")
// 				{
// 					// List all customers
// 					customer.GET("/list", customerHandler.ListCustomers)

// 					// Show the form to create a new customer
// 					customer.GET("/create", customerHandler.CreateCustomerForm)

// 					// Process the form submission to create a new customer
// 					customer.POST("/create", customerHandler.CreateCustomer)

// 					// Show the form to update a customer
// 					customer.GET("/update", customerHandler.UpdateCustomerForm)

// 					// Process the form submission to update a customer
// 					customer.POST("/update", customerHandler.UpdateCustomer)

// 					// Delete a customer
// 					customer.GET("/delete", customerHandler.DeleteCustomer)
// 				}
// 			}

// 			return r
// 		},
// 	}

// 	for _, provider := range providers {
// 		if err := container.Provide(provider); err != nil {
// 			log.Critical("Failed to provide dependency", log.Fstring("provider", fmt.Sprintf("%T", provider)))
// 			return nil, err
// 		}
// 	}

// 	log.Info("Container built successfully")
// 	return container, nil
// }
