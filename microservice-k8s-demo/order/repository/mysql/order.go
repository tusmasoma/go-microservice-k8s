package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (or *orderRepository) Get(ctx context.Context, id string) (*entity.Order, error) {
	query := `
	SELECT
		Orders.id,
		Orders.customer_id,
		Orders.order_date,
		OrderLines.catalog_item_id,
		OrderLines.count
	FROM
   		Orders
	INNER JOIN
    	OrderLines ON Orders.id = OrderLines.order_id
	WHERE
		Orders.id = ?
	`

	rows, err := or.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order entity.Order
	for rows.Next() {
		var orderLine entity.OrderLine
		if err = rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.OrderDate,
			&orderLine.CatalogItemID,
			&orderLine.Count,
		); err != nil {
			return nil, err
		}
		order.OrderLines = append(order.OrderLines, orderLine)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (or *orderRepository) List(ctx context.Context) ([]entity.Order, error) {
	query := `
	SELECT
		Orders.id,
		Orders.customer_id,
		Orders.order_date,
		OrderLines.catalog_item_id,
		OrderLines.count
	FROM
   		Orders
	INNER JOIN
    	OrderLines ON Orders.id = OrderLines.order_id
	`

	rows, err := or.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	orderMap := make(map[string]*entity.Order)

	for rows.Next() {
		var id string
		var orderLine entity.OrderLine
		var orderDate time.Time
		var customerID string

		if err = rows.Scan(&id, &customerID, &orderDate, &orderLine.CatalogItemID, &orderLine.Count); err != nil {
			return nil, err
		}

		order, exists := orderMap[id]
		if !exists {
			order = &entity.Order{
				ID:         id,
				CustomerID: customerID,
				OrderDate:  orderDate,
			}
			orderMap[id] = order
			orders = append(orders, *order)
		}

		order.OrderLines = append(order.OrderLines, orderLine)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *orderRepository) Create(ctx context.Context, order entity.Order) error {
	tx, err := or.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
			panic(p)      // re-throw the panic after Rollback
		} else if err != nil {
			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
		} else {
			err = tx.Commit() // The error is checked here.
		}
	}()

	query := `
	INSERT INTO Orders (id, customer_id, order_date)
	VALUES (?, ?, ?)
	`
	if _, err = tx.ExecContext(
		ctx,
		query,
		order.ID,
		order.CustomerID,
		order.OrderDate,
	); err != nil {
		return err
	}

	query = `
	INSERT INTO OrderLines (order_id, catalog_item_id, count) VALUES`
	values := make([]interface{}, 0, len(order.OrderLines)*3) //nolint:gomnd // 3 is the number of columns.
	for i, line := range order.OrderLines {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?)"
		values = append(values, order.ID, line.CatalogItemID, line.Count)
	}

	if _, err = tx.ExecContext(ctx, query, values...); err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) Delete(ctx context.Context, id string) error {
	// Application-level responsibility:
	// This method is responsible for deleting both the order and its associated order lines.
	// Although the database could handle this automatically with ON DELETE CASCADE,
	// we are managing the deletion process here at the application level for greater flexibility.
	tx, err := or.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
			panic(p)      // re-throw the panic after Rollback
		} else if err != nil {
			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
		} else {
			err = tx.Commit() // The error is checked here.
		}
	}()

	query := `
	DELETE FROM OrderLines WHERE order_id = ?
	`
	if _, err = tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	query = `
	DELETE FROM Orders WHERE id = ?
	`
	if _, err = tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
