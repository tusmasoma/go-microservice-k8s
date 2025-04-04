package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Order struct {
	order.Order
}

func (o *Order) Proto() *order.Order {
	if o == nil {
		return nil
	}
	return &o.Order
}

func (o *Order) GetOrderLines() OrderLines {
	if o == nil || o.Order.OrderLines == nil {
		return nil
	}
	orderLines := make(OrderLines, len(o.Order.GetOrderLines()))
	for i, orderLine := range o.Order.GetOrderLines() {
		orderLines[i] = &OrderLine{OrderLine: orderLine}
	}
	return orderLines
}

type Orders []*Order

func (os Orders) Proto() []*order.Order {
	orders := make([]*order.Order, len(os))
	for i, order := range os {
		orders[i] = order.Proto()
	}
	return orders
}

type OrderLine struct {
	*order.OrderLine
}

func (ol *OrderLine) Proto() *order.OrderLine {
	if ol == nil {
		return nil
	}
	return ol.OrderLine
}

type OrderLines []*OrderLine

func (ols OrderLines) Proto() []*order.OrderLine {
	orderLines := make([]*order.OrderLine, len(ols))
	for i, orderLine := range ols {
		orderLines[i] = orderLine.Proto()
	}
	return orderLines
}

func NewOrder(id, customerID string, orderDate time.Time, orderLines OrderLines) (*Order, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if customerID == "" {
		return nil, errors.New("customerID is required")
	}
	if len(orderLines) == 0 {
		return nil, errors.New("order lines is required")
	}
	return &Order{
		Order: order.Order{
			Id:         id,
			CustomerId: customerID,
			OrderDate:  timestamppb.New(orderDate),
			OrderLines: orderLines.Proto(),
		},
	}, nil
}

func CreateOrder(customerID string, orderLines OrderLines) (*Order, error) {
	if customerID == "" {
		return nil, errors.New("customerID is required")
	}
	if len(orderLines) == 0 {
		return nil, errors.New("order lines is required")
	}
	return &Order{
		Order: order.Order{
			Id:         uuid.NewString(),
			CustomerId: customerID,
			OrderDate:  timestamppb.New(time.Now()),
			OrderLines: orderLines.Proto(),
		},
	}, nil
}
