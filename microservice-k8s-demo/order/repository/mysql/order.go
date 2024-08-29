package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/repository"
)

// As an alternative approach, both the Order and OrderLine tables could be managed within separate repositories,
// with CRUD operations for each table implemented independently. In this scenario, the service layer (or another coordinating layer)
// would be responsible for ensuring consistency across operations that affect both tables.
// This approach can improve modularity and make the codebase easier to maintain, but requires careful coordination at the service level.

// In this case, we are implementing both the Order and OrderLine tables within the same repository.
// This approach simplifies the management of related data and ensures consistency within the repository itself.
// While this can reduce the complexity at the service layer, it may result in a larger, more tightly coupled repository.

type orderRepository struct {
	db SQLExecutor
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (or *orderRepository) Get(ctx context.Context, id string) (*entity.OrderModel, error) {
	executor := or.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, customer_id, date
	FROM Orders
	WHERE id = ?
	LIMIT 1
	`

	row := executor.QueryRowContext(ctx, query, id)

	var order entity.OrderModel
	if err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.OrderDate,
	); err != nil {
		return nil, err
	}

	return &order, nil
}

func (or *orderRepository) List(ctx context.Context) ([]entity.OrderModel, error) {
	executor := or.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, customer_id, date
	FROM Orders
	`

	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.OrderModel
	for rows.Next() {
		var order entity.OrderModel
		if err = rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.OrderDate,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *orderRepository) Create(ctx context.Context, order entity.OrderModel) error {
	executor := or.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	INSERT INTO Orders (id, customer_id, date)
	VALUES (?, ?, ?)
	`
	_, err := executor.ExecContext(
		ctx,
		query,
		order.ID,
		order.CustomerID,
		order.OrderDate,
	)
	return err
}

func (or *orderRepository) Delete(ctx context.Context, id string) error {
	executor := or.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	DELETE FROM Orders WHERE id = ?
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}

type orderLineRepository struct {
	db SQLExecutor
}

func NewOrderLineRepository(db *sql.DB) repository.OrderLineRepository {
	return &orderLineRepository{
		db: db,
	}
}

func (olr *orderLineRepository) List(ctx context.Context, orderID string) ([]entity.OrderLineModel, error) {
	executor := olr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT order_id, catalog_item_id, count
	FROM OrderLines
	WHERE order_id = ?
	`

	rows, err := executor.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderLines []entity.OrderLineModel
	for rows.Next() {
		var orderLine entity.OrderLineModel
		if err = rows.Scan(
			&orderLine.OrderID,
			&orderLine.CatalogItemID,
			&orderLine.Count,
		); err != nil {
			return nil, err
		}
		orderLines = append(orderLines, orderLine)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orderLines, nil
}

func (olr *orderLineRepository) Create(ctx context.Context, orderLine entity.OrderLineModel) error {
	executor := olr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	INSERT INTO OrderLines (order_id, catalog_item_id, count)
	VALUES (?, ?, ?)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		orderLine.OrderID,
		orderLine.CatalogItemID,
		orderLine.Count,
	); err != nil {
		return err
	}
	return nil
}

func (olr *orderLineRepository) BatchCreate(ctx context.Context, orderLines []entity.OrderLineModel) error {
	executor := olr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	INSERT INTO OrderLines (order_id, catalog_item_id, count) VALUES
	`

	values := make([]interface{}, 0, len(orderLines)*3) //nolint:gomnd // 3 is the number of columns.
	for i, line := range orderLines {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?)"
		values = append(values, line.OrderID, line.CatalogItemID, line.Count)
	}

	if _, err := executor.ExecContext(ctx, query, values...); err != nil {
		return err
	}

	return nil
}

func (olr *orderLineRepository) Delete(ctx context.Context, orderID, catalogItemID string) error {
	executor := olr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	DELETE FROM OrderLines WHERE order_id = ? AND catalog_item_id = ?
	`

	if _, err := executor.ExecContext(ctx, query, orderID); err != nil {
		return err
	}
	return nil
}

func (olr *orderLineRepository) BatchDelete(ctx context.Context, orderID string) error {
	executor := olr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	DELETE FROM OrderLines WHERE order_id = ?
	`

	if _, err := executor.ExecContext(ctx, query, orderID); err != nil {
		return err
	}
	return nil
}

// func (or *orderRepository) Get(ctx context.Context, id string) (*entity.Order, error) {
// 	query := `
// 	SELECT
// 		Orders.id,
// 		Orders.customer_id,
// 		Orders.order_date,
// 		OrderLines.catalog_item_id,
// 		OrderLines.count
// 	FROM
//    		Orders
// 	INNER JOIN
//     	OrderLines ON Orders.id = OrderLines.order_id
// 	WHERE
// 		Orders.id = ?
// 	`

// 	rows, err := or.db.QueryContext(ctx, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var order entity.Order
// 	for rows.Next() {
// 		var orderLine entity.OrderLine
// 		if err = rows.Scan(
// 			&order.ID,
// 			&order.CustomerID,
// 			&order.OrderDate,
// 			&orderLine.CatalogItemID,
// 			&orderLine.Count,
// 		); err != nil {
// 			return nil, err
// 		}
// 		order.OrderLines = append(order.OrderLines, orderLine)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return &order, nil
// }

// func (or *orderRepository) List(ctx context.Context) ([]entity.Order, error) {
// 	query := `
// 	SELECT
// 		Orders.id,
// 		Orders.customer_id,
// 		Orders.order_date,
// 		OrderLines.catalog_item_id,
// 		OrderLines.count
// 	FROM
//    		Orders
// 	INNER JOIN
//     	OrderLines ON Orders.id = OrderLines.order_id
// 	`

// 	rows, err := or.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var orders []entity.Order
// 	orderMap := make(map[string]*entity.Order)

// 	for rows.Next() {
// 		var id string
// 		var orderLine entity.OrderLine
// 		var orderDate time.Time
// 		var customerID string

// 		if err = rows.Scan(&id, &customerID, &orderDate, &orderLine.CatalogItemID, &orderLine.Count); err != nil {
// 			return nil, err
// 		}

// 		order, exists := orderMap[id]
// 		if !exists {
// 			order = &entity.Order{
// 				ID:         id,
// 				CustomerID: customerID,
// 				OrderDate:  orderDate,
// 			}
// 			orderMap[id] = order
// 			orders = append(orders, *order)
// 		}

// 		order.OrderLines = append(order.OrderLines, orderLine)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

// func (or *orderRepository) Create(ctx context.Context, order entity.Order) error {
// 	tx, err := or.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
// 	if err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
// 			panic(p)      // re-throw the panic after Rollback
// 		} else if err != nil {
// 			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
// 		} else {
// 			err = tx.Commit() // The error is checked here.
// 		}
// 	}()

// 	query := `
// 	INSERT INTO Orders (id, customer_id, order_date)
// 	VALUES (?, ?, ?)
// 	`
// 	if _, err = tx.ExecContext(
// 		ctx,
// 		query,
// 		order.ID,
// 		order.CustomerID,
// 		order.OrderDate,
// 	); err != nil {
// 		return err
// 	}

// 	query = `
// 	INSERT INTO OrderLines (order_id, catalog_item_id, count) VALUES`
// 	values := make([]interface{}, 0, len(order.OrderLines)*3) //nolint:gomnd // 3 is the number of columns.
// 	for i, line := range order.OrderLines {
// 		if i > 0 {
// 			query += ", "
// 		}
// 		query += "(?, ?, ?)"
// 		values = append(values, order.ID, line.CatalogItemID, line.Count)
// 	}

// 	if _, err = tx.ExecContext(ctx, query, values...); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (or *orderRepository) Delete(ctx context.Context, id string) error {
// 	// Application-level responsibility:
// 	// This method is responsible for deleting both the order and its associated order lines.
// 	// Although the database could handle this automatically with ON DELETE CASCADE,
// 	// we are managing the deletion process here at the application level for greater flexibility.
// 	tx, err := or.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
// 	if err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
// 			panic(p)      // re-throw the panic after Rollback
// 		} else if err != nil {
// 			tx.Rollback() //nolint:errcheck // The error is checked in the outer function.
// 		} else {
// 			err = tx.Commit() // The error is checked here.
// 		}
// 	}()

// 	query := `
// 	DELETE FROM OrderLines WHERE order_id = ?
// 	`
// 	if _, err = tx.ExecContext(ctx, query, id); err != nil {
// 		return err
// 	}

// 	query = `
// 	DELETE FROM Orders WHERE id = ?
// 	`
// 	if _, err = tx.ExecContext(ctx, query, id); err != nil {
// 		return err
// 	}

// 	return nil
// }
