package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// domain model
type Order struct {
	ID         string       `json:"id"`
	CustomerID string       `json:"customer_id"`
	OrderDate  *time.Time   `json:"order_date"`
	OrderLines []*OrderLine `json:"order_lines"`
	TotalPrice float64      `json:"total_price"`
}

type OrderLine struct {
	Count         int    `json:"count"`
	CatalogItemID string `json:"catalog_item_id"`
}

func NewOrder(id, customerID string, orderDate *time.Time, orderLines []*OrderLine) (*Order, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if customerID == "" {
		return nil, errors.New("customerID is required")
	}
	if orderDate == nil {
		orderDate = new(time.Time)
	}
	order := &Order{
		ID:         id,
		CustomerID: customerID,
		OrderDate:  orderDate,
		OrderLines: orderLines,
	}

	// order.TotalPrice = order.GetTotalPrice()
	return order, nil
}

func NewOrderLine(count int, itemID string) (*OrderLine, error) {
	if count <= 0 {
		return nil, errors.New("count must be greater than 0")
	}
	if itemID == "" {
		return nil, errors.New("catalogItemID is required")
	}
	return &OrderLine{
		Count:         count,
		CatalogItemID: itemID,
	}, nil
}

// func (o *Order) GetTotalPrice() float64 {
// 	var total float64
// 	for _, ol := range o.OrderLines {
// 		total += ol.CatalogItem.Price * float64(ol.Count)
// 	}
// 	return total
// }
