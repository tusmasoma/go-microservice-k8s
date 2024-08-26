package entity

import "time"

type Order struct {
	ID         string      `json:"id" db:"id"`
	CustomerID string      `json:"customer_id" db:"customer_id"`
	OrderDate  time.Time   `json:"order_date" db:"date"`
	OrderLines []OrderLine `json:"order_lines"`
}

type OrderLine struct {
	Count         int    `json:"count" db:"count"`
	CatalogItemID string `json:"catalog_item_id" db:"catalog_item_id"`
}

// NewOrder creates a new order.
