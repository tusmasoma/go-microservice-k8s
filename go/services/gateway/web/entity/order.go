package entity

import (
	"github.com/tusmasoma/go-microservice-k8s/proto/catalog"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
	"github.com/tusmasoma/go-microservice-k8s/proto/order"
	"github.com/tusmasoma/go-microservice-k8s/proto/web"
)

type Order struct {
	web.Order
}

func (o *Order) Proto() *web.Order {
	if o == nil {
		return nil
	}
	return &o.Order
}

func NewOrder(order *order.Order, customer *customer.Customer, catalogItem []*catalog.CatalogItem) *Order {
	return &Order{
		Order: web.Order{
			Id:         order.GetId(),
			Customer:   NewCustomer(customer).Proto(),
			OrderDate:  order.GetOrderDate(),
			OrderLines: newOrderLines(order.GetOrderLines(), catalogItem).Proto(),
		},
	}
}

type Orders []*Order

func (os Orders) Proto() []*web.Order {
	orders := make([]*web.Order, len(os))
	for i, order := range os {
		orders[i] = order.Proto()
	}
	return orders
}

type OrderLine struct {
	*web.OrderLine
}

func newOrderLine(line *order.OrderLine, catalogItem *catalog.CatalogItem) *OrderLine {
	if line.GetCatalogItemId() != catalogItem.GetId() {
		return nil
	}
	return &OrderLine{
		OrderLine: &web.OrderLine{
			Count: line.GetCount(),
			CatalogItem: &web.CatalogItem{
				Id:    catalogItem.GetId(),
				Name:  catalogItem.GetName(),
				Price: catalogItem.GetPrice(),
			},
		},
	}
}

func (ol *OrderLine) Proto() *web.OrderLine {
	if ol == nil {
		return nil
	}
	return ol.OrderLine
}

type OrderLines []*OrderLine

func newOrderLines(lines []*order.OrderLine, catalogItems []*catalog.CatalogItem) OrderLines {
	itemMap := make(map[string]*catalog.CatalogItem, len(catalogItems))
	for _, item := range catalogItems {
		itemMap[item.GetId()] = item
	}
	var result OrderLines
	for _, line := range lines {
		catalogItem := itemMap[line.GetCatalogItemId()]
		if catalogItem == nil {
			continue
		}
		orderLine := newOrderLine(line, catalogItem)
		if orderLine != nil {
			result = append(result, orderLine)
		}
	}
	return result
}

func (ols OrderLines) Proto() []*web.OrderLine {
	lines := make([]*web.OrderLine, len(ols))
	for i, line := range ols {
		lines[i] = line.Proto()
	}
	return lines
}
