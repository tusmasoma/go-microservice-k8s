package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/config"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/gateway"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/repository/mysql"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/usecase"

	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/proto"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found", log.Ferror(err))
	}

	var addr string
	flag.StringVar(&addr, "addr", ":8081", "tcp host:port to connect")
	flag.Parse()

	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	container, err := BuildContainer(mainCtx)
	if err != nil {
		log.Critical("Failed to build container", log.Ferror(err))
		return
	}

	err = container.Invoke(func(grpcHandler pb.CustomerServiceServer, config *config.ServerConfig) {
		lis, err := net.Listen("tcp", addr) //nolint:govet // This is not a mistake
		if err != nil {
			log.Critical("Failed to listen", log.Ferror(err))
		}

		srv := grpc.NewServer()

		pb.RegisterCustomerServiceServer(srv, grpcHandler)

		reflection.Register(srv)

		log.Info("Server started", log.Fstring("addr", addr))

		go func() {
			if err = srv.Serve(lis); err != nil {
				log.Critical("Failed to serve", log.Ferror(err))
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		log.Info("Server stopping...")
		srv.GracefulStop()
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
		mysql.NewCustomerRepository,
		usecase.NewCustomerUsecase,
		gateway.NewCustomerHandler,
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
