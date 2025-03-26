package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tusmasoma/go-microservice-k8s/pkg/config"
	"github.com/tusmasoma/go-microservice-k8s/services/gateway/web"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found", log.Ferror(err))
	}
	var addr string
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()
	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()
	// --- DI ---
	serverConfig, err := config.NewServerConfig(mainCtx)
	if err != nil {
		log.Critical("Failed to initialize server config", log.Ferror(err))
		return
	}
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"https://*", "http://*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:     []string{"Link", "Authorization"},
		AllowCredentials:   true,
		MaxAge:             serverConfig.PreflightCacheDurationSec,
		OptionsPassthrough: false,
	}))
	params, err := setupHandlerParams(mainCtx)
	if err != nil {
		return
	}
	web.NewWeb(params).Register(mux)
	// --- Server Run ---
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}
	log.Info("Server running...")
	// --- Graceful shutdown ---
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
	tctx, cancelShutdown := context.WithTimeout(context.Background(), serverConfig.GracefulShutdownTimeout)
	defer cancelShutdown()
	if err = srv.Shutdown(tctx); err != nil {
		log.Error("Failed to shutdown http server", log.Ferror(err))
	}
	log.Info("Server exited")
}
