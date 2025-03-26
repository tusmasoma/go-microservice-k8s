package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/tusmasoma/go-microservice-k8s/go/services/gateway/web"
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	"google.golang.org/grpc"
)

func setupHandlerParams(ctx context.Context) (*web.HandlerParams, error) {
	customerClient, err := newCustomerClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	catalogClient, err := newCatalogClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	orderClient, err := newOrderClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &web.HandlerParams{
		CustomerClient: customerClient,
		CatalogClient:  catalogClient,
		OrderClient:    orderClient,
	}, nil
}

func newCustomerClient(ctx context.Context, opts ...grpc.DialOption) (customer.CustomerServiceClient, error) {
	return newClientWithDial(ctx, customer.NewCustomerServiceClient, opts...)
}

func newCatalogClient(ctx context.Context, opts ...grpc.DialOption) (catalog.CatalogServiceClient, error) {
	return newClientWithDial(ctx, catalog.NewCatalogServiceClient, opts...)
}

func newOrderClient(ctx context.Context, opts ...grpc.DialOption) (order.OrderServiceClient, error) {
	return newClientWithDial(ctx, order.NewOrderServiceClient, opts...)
}

const defaultPort = 9000

func newClientWithDial[T any](_ context.Context, newClient func(grpc.ClientConnInterface) T, opts ...grpc.DialOption) (T, error) {
	serviceName := serviceNameFromType[T]()
	target := fmt.Sprintf("%s-server:%d", serviceName, defaultPort)
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return *new(T), fmt.Errorf("grpc: failed to dial %s: %w", target, err)
	}
	return newClient(conn), nil
}

func serviceNameFromType[T any]() string {
	t := reflect.TypeFor[T]()
	typeName := t.String()
	packageName := strings.Split(typeName, ".")[0]
	return packageName
}
