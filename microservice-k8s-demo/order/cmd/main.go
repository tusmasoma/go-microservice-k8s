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

	catalog_pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/catalog/proto"
	cusotmer_pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/proto"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/config"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/gateway"
	pb "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/proto"
	catalogservice "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository/catalog_service"
	customerservice "github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository/customer_service"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository/mysql"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/usecase"
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

	container, err := BuildContainer(mainCtx)
	if err != nil {
		log.Critical("Failed to build container", log.Ferror(err))
		return
	}

	err = container.Invoke(func(grpcHandler pb.OrderServiceServer, config *config.ServerConfig) {
		lis, err := net.Listen("tcp", addr) //nolint:govet // This is not a mistake
		if err != nil {
			log.Critical("Failed to listen", log.Ferror(err))
		}

		srv := grpc.NewServer()

		pb.RegisterOrderServiceServer(srv, grpcHandler)

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
		mysql.NewOrderRepository,
		NewCustomerServiceClient,
		NewCatalogServiceClient,
		customerservice.NewCustomerRepository,
		catalogservice.NewCatalogItemRepository,
		usecase.NewOrderUseCase,
		gateway.NewOrderHandler,
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

func NewCatalogServiceClient() catalog_pb.CatalogServiceClient {
	conn, _ := grpc.Dial("catalog-service:8082", grpc.WithInsecure()) //nolint:staticcheck // ignore deprecation
	return catalog_pb.NewCatalogServiceClient(conn)
}

func NewCustomerServiceClient() cusotmer_pb.CustomerServiceClient {
	conn, _ := grpc.Dial("customer-service:8081", grpc.WithInsecure()) //nolint:staticcheck // ignore deprecation
	return cusotmer_pb.NewCustomerServiceClient(conn)
}
