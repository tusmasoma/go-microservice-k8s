package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/config"
	"github.com/tusmasoma/microservice-k8s-demo/catalog/gateway/web/handler"
	"github.com/tusmasoma/microservice-k8s-demo/catalog/repository/mysql"

	"go.uber.org/dig"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/usecase"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found", log.Ferror(err))
	}

	var addr string
	flag.StringVar(&addr, "addr", ":8082", "tcp host:port to connect")
	flag.Parse()

	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	// This is framework-agnostic and can be replaced with any HTTP framework like net/http, Gin, or Echo
	// The following example uses net/http.
	container, err := BuildContainer(mainCtx)
	if err != nil {
		log.Critical("Failed to build container", log.Ferror(err))
		return
	}

	err = container.Invoke(func(router *chi.Mux, config *config.ServerConfig) {
		srv := &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
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

		tctx, cancelShutdown := context.WithTimeout(context.Background(), config.GracefulShutdownTimeout)
		defer cancelShutdown()

		if err = srv.Shutdown(tctx); err != nil {
			log.Error("Failed to shutdown http server", log.Ferror(err))
		}
		log.Info("Server exited")
	})
	if err != nil {
		log.Critical("Failed to start server", log.Ferror(err))
		return
	}
}

func BuildContainer(ctx context.Context) (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		log.Error("Failed to provide context")
		return nil, err
	}

	providers := []interface{}{
		config.NewServerConfig,
		config.NewDBConfig,
		mysql.NewMySQLDB,
		mysql.NewTransactionRepository,
		mysql.NewCatalogItemRepository,
		usecase.NewCatalogItemUseCase,
		handler.NewCatalogItemHandler,
		func(
			serverConfig *config.ServerConfig,
			catalogHandler handler.CatalogItemHandler,
		) *chi.Mux {
			r := chi.NewRouter()
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:     []string{"https://*", "http://*"},
				AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
				ExposedHeaders:     []string{"Link", "Authorization"},
				AllowCredentials:   true,
				MaxAge:             serverConfig.PreflightCacheDurationSec,
				OptionsPassthrough: false,
			}))

			r.Route("/catalog", func(r chi.Router) {
				r.Get("/list", catalogHandler.ListCatalogItems)
			})

			return r
		},
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			log.Critical("Failed to provide dependency", log.Fstring("provider", fmt.Sprintf("%T", provider)))
			return nil, err
		}
	}

	log.Info("Container built successfully")
	return container, nil
}
