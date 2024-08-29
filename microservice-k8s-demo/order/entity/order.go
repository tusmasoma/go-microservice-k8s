package entity

import (
	"time"
)

// domain model
type Order struct {
	ID         string      `json:"id" db:"id"`
	Customer   Customer    `json:"customer"`
	OrderDate  time.Time   `json:"order_date" db:"date"`
	OrderLines []OrderLine `json:"order_lines"`
}

type OrderLine struct {
	Count       int         `json:"count" db:"count"`
	CatalogItem CatalogItem `json:"catalog_item"`
}

// data model
type OrderModel struct {
	ID         string    `json:"id" db:"id"`
	CustomerID string    `json:"customer_id" db:"customer_id"`
	OrderDate  time.Time `json:"order_date" db:"date"`
}

type OrderLineModel struct {
	OrderID       string `json:"order_id" db:"order_id"`
	CatalogItemID string `json:"catalog_item_id" db:"catalog_item_id"`
	Count         int    `json:"count" db:"count"`
}

// func NewOrder(customerID string, orderLines []OrderLine) (*Order, error) {
// 	if customerID == "" {
// 		return nil, errors.New("customerID is required")
// 	}
// 	if len(orderLines) == 0 {
// 		return nil, errors.New("orderLines is required")
// 	}
// 	for _, ol := range orderLines {
// 		if ol.Count <= 0 {
// 			return nil, errors.New("count must be greater than 0")
// 		}
// 		if ol.CatalogItemID == "" {
// 			return nil, errors.New("catalogItemID is required")
// 		}
// 	}
// 	return &Order{
// 		ID:         uuid.New().String(),
// 		CustomerID: customerID,
// 		OrderDate:  time.Now(),
// 		OrderLines: orderLines,
// 	}, nil
// }

func (o *Order) TotalPrice() float64 {
	var total float64 = 0
	for _, ol := range o.OrderLines {
		total += ol.CatalogItem.Price * float64(ol.Count)
	}
	return total
}
