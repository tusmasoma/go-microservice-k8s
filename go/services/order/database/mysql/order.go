package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/tusmasoma/go-microservice-k8s/go/pkg/log"
	"github.com/tusmasoma/go-microservice-k8s/go/services/order/database"
	"github.com/tusmasoma/go-microservice-k8s/go/services/order/entity"
	pb "github.com/tusmasoma/go-microservice-k8s/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type order struct {
	db *sql.DB
}

func NewOrder(db *sql.DB) database.Order {
	return &order{
		db: db,
	}
}

func (or *order) Get(ctx context.Context, id string) (*entity.Order, error) {
	query := `
	SELECT id, customer_id, order_date
	FROM Orders
	WHERE id = ?
	LIMIT 1
	`
	row := or.db.QueryRowContext(ctx, query, id)
	var order entity.Order
	var orderDate time.Time
	if err := row.Scan(
		&order.Id,
		&order.CustomerId,
		&orderDate, // sql not support timestamppb.Timestamp
	); err != nil {
		return nil, err
	}
	order.OrderDate = timestamppb.New(orderDate)
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
	var lines entity.OrderLines
	for rows.Next() {
		var line entity.OrderLine
		line.OrderLine = &pb.OrderLine{} // initialize here to prevent nil pointers
		if err = rows.Scan(
			&line.CatalogItemId,
			&line.Count,
		); err != nil {
			return nil, err
		}
		lines = append(lines, &line)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	order.OrderLines = lines.Proto()
	return &order, nil
}

func (or *order) List(ctx context.Context) (entity.Orders, error) {
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
	var orders entity.Orders
	orderMap := make(map[string]*entity.Order)
	for rows.Next() {
		var (
			orderID       string
			customerID    string
			orderDate     time.Time
			catalogItemID string
			count         int64
		)
		if err = rows.Scan(
			&orderID,
			&customerID,
			&orderDate,
			&catalogItemID,
			&count,
		); err != nil {
			return nil, err
		}
		order, exists := orderMap[orderID]
		if !exists {
			order = &entity.Order{
				Order: pb.Order{
					Id:         orderID,
					CustomerId: customerID,
					OrderDate:  timestamppb.New(orderDate),
				},
			}
			orderMap[orderID] = order
			orders = append(orders, order)
		}
		orderLine := entity.OrderLine{
			OrderLine: &pb.OrderLine{
				CatalogItemId: catalogItemID,
				Count:         count,
			},
		}
		order.OrderLines = append(order.OrderLines, orderLine.Proto())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (or *order) Create(ctx context.Context, order *entity.Order) error {
	if err := or.transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err := or.createOrder(ctx, tx, order); err != nil {
			return err
		}
		if err := or.createOrderLines(ctx, tx, order.GetId(), order.GetOrderLines()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (or *order) createOrder(ctx context.Context, tx *sql.Tx, order *entity.Order) error {
	query := `
	INSERT INTO Orders (id, customer_id, order_date)
	VALUES (?, ?, ?)
	`
	if _, err := tx.ExecContext(
		ctx,
		query,
		order.GetId(),
		order.GetCustomerId(),
		order.GetOrderDate().AsTime(),
	); err != nil {
		return err
	}
	return nil
}

func (or *order) createOrderLines(ctx context.Context, tx *sql.Tx, orderID string, lines entity.OrderLines) error {
	query := `
	INSERT INTO OrderLines (order_id, catalog_item_id, count) VALUES`
	values := make([]interface{}, 0, len(lines)*3) //nolint:mnd // 3 is the number of columns.
	for i, line := range lines {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?)"
		values = append(values, orderID, line.GetCatalogItemId(), line.GetCount())
	}
	if _, err := tx.ExecContext(ctx, query, values...); err != nil {
		return err
	}
	return nil
}

func (or *order) Delete(ctx context.Context, id string) error {
	if err := or.transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err := or.batchDeleteOrderLines(ctx, tx, id); err != nil {
			return err
		}
		if err := or.deleteOrders(ctx, tx, id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (or *order) batchDeleteOrderLines(ctx context.Context, tx *sql.Tx, orderID string) error {
	query := "DELETE FROM OrderLines WHERE order_id = ?"
	if _, err := tx.ExecContext(ctx, query, orderID); err != nil {
		return err
	}
	return nil
}

func (or *order) deleteOrders(ctx context.Context, tx *sql.Tx, orderID string) error {
	query := "DELETE FROM Orders WHERE id = ?"
	if _, err := tx.ExecContext(ctx, query, orderID); err != nil {
		return err
	}
	return nil
}

func (or *order) transaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := or.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()
	if err = fn(ctx, tx); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
