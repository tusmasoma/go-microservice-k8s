package mysql

import (
	"context"
	"database/sql"
	"time"

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

type orderModel struct {
	ID         string    `db:"id"`
	CustomerID string    `db:"customer_id"`
	OrderDate  time.Time `db:"order_date"`
}

type orderLineModel struct {
	OrderID       string `db:"order_id"`
	CatalogItemID string `db:"catalog_item_id"`
	Count         int    `db:"count"`
}

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

	var om orderModel
	if err := row.Scan(
		&om.ID,
		&om.CustomerID,
		&om.OrderDate,
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

	var olms []orderLineModel
	for rows.Next() {
		var olm orderLineModel
		if err = rows.Scan(
			&olm.CatalogItemID,
			&olm.Count,
		); err != nil {
			return nil, err
		}
		olms = append(olms, olm)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Mapping to entity.Order and entity.OrderLine
	orderLines := make([]*entity.OrderLine, 0, len(olms))
	for _, line := range olms {
		orderLine, err := entity.NewOrderLine(line.Count, line.CatalogItemID) //nolint:govet // err shadowed
		if err != nil {
			return nil, err
		}
		orderLines = append(orderLines, orderLine)
	}

	order, err := entity.NewOrder(om.ID, om.CustomerID, &om.OrderDate, orderLines)
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
		var om orderModel
		var olm orderLineModel

		if err = rows.Scan(
			&om.ID,
			&om.CustomerID,
			&om.OrderDate,
			&olm.CatalogItemID,
			&olm.Count,
		); err != nil {
			return nil, err
		}

		order, exists := orderMap[om.ID]
		if !exists {
			order, err = entity.NewOrder(om.ID, om.CustomerID, &om.OrderDate, nil)
			if err != nil {
				return nil, err
			}
			orderMap[om.ID] = order
			orders = append(orders, order)
		}

		orderLine, err := entity.NewOrderLine(olm.Count, olm.CatalogItemID) //nolint:govet // err shadowed
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

	om := orderModel{
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
		om.ID,
		om.CustomerID,
		om.OrderDate,
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

		olm := orderLineModel{
			OrderID:       order.ID,
			CatalogItemID: line.CatalogItemID,
			Count:         line.Count,
		}
		values = append(values, olm.OrderID, olm.CatalogItemID, olm.Count)
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
