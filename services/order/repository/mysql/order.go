package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-microservice-k8s/services/order/entity"
	"github.com/tusmasoma/go-microservice-k8s/services/order/repository"
)

// As an alternative approach, both the Order and OrderLine tables could be managed within separate repositories,
// with CRUD operations for each table implemented independently. In this scenario, the service layer (or another coordinating layer)
// would be responsible for ensuring consistency across operations that affect both tables.
// This approach can improve modularity and make the codebase easier to maintain, but requires careful coordination at the service level.

// In this case, we are implementing both the Order and OrderLine tables within the same repository.
// This approach simplifies the management of related data and ensures consistency within the repository itself.
// While this can reduce the complexity at the service layer, it may result in a larger, more tightly coupled repository.

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (or *orderRepository) Get(ctx context.Context, id string) (*entity.Order, error) {
	// Orders table query
	query := `
	SELECT id, customer_id, order_date
	FROM Orders
	WHERE id = ?
	LIMIT 1
	`

	row := or.db.QueryRowContext(ctx, query, id)

	var orderModel entity.OrderModel
	if err := row.Scan(
		&orderModel.ID,
		&orderModel.CustomerID,
		&orderModel.OrderDate,
	); err != nil {
		return nil, err
	}

	// OrderLines table query
	query = `
	SELECT catalog_item_id, count
	FROM OrderLines
	WHERE order_id = ?
	`

	rows, err := or.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderModelLines []entity.OrderLineModel
	for rows.Next() {
		var orderLineModel entity.OrderLineModel
		if err = rows.Scan(
			&orderLineModel.CatalogItemID,
			&orderLineModel.Count,
		); err != nil {
			return nil, err
		}
		orderModelLines = append(orderModelLines, orderLineModel)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Mapping to entity.Order and entity.OrderLine
	orderLines := make([]*entity.OrderLine, 0, len(orderModelLines))
	for _, line := range orderModelLines {
		orderLine, err := entity.NewOrderLine(line.Count, line.CatalogItemID)
		if err != nil {
			return nil, err
		}
		orderLines = append(orderLines, orderLine)
	}

	order, err := entity.NewOrder(orderModel.ID, orderModel.CustomerID, &orderModel.OrderDate, orderLines)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (or *orderRepository) List(ctx context.Context) ([]*entity.Order, error) {
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

	var orders []*entity.Order
	orderMap := make(map[string]*entity.Order)

	for rows.Next() {
		var orderModel entity.OrderModel
		var orderLineModel entity.OrderLineModel

		if err = rows.Scan(
			&orderModel.ID,
			&orderModel.CustomerID,
			&orderModel.OrderDate,
			&orderLineModel.CatalogItemID,
			&orderLineModel.Count,
		); err != nil {
			return nil, err
		}

		order, exists := orderMap[orderModel.ID]
		if !exists {
			order, err := entity.NewOrder(orderModel.ID, orderModel.CustomerID, &orderModel.OrderDate, nil)
			if err != nil {
				return nil, err
			}
			orderMap[orderModel.ID] = order
			orders = append(orders, order)
		}

		orderLine, err := entity.NewOrderLine(orderLineModel.Count, orderLineModel.CatalogItemID)
		if err != nil {
			return nil, err
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

	orderModel := entity.OrderModel{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		OrderDate:  *order.OrderDate,
	}

	query := `
	INSERT INTO Orders (id, customer_id, order_date)
	VALUES (?, ?, ?)
	`
	if _, err = tx.ExecContext(
		ctx,
		query,
		orderModel.ID,
		orderModel.CustomerID,
		orderModel.OrderDate,
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

		orderLineModel := entity.OrderLineModel{
			OrderID:       order.ID,
			CatalogItemID: line.CatalogItemID,
			Count:         line.Count,
		}
		values = append(values, orderLineModel.OrderID, orderLineModel.CatalogItemID, orderLineModel.Count)
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
