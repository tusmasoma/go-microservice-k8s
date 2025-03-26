package mysql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/tusmasoma/go-microservice-k8s/services/catalog/database"
	"github.com/tusmasoma/go-microservice-k8s/services/catalog/entity"
)

type catalogItem struct {
	db *sql.DB
}

func NewCatalogItem(db *sql.DB) database.CatalogItem {
	return &catalogItem{
		db: db,
	}
}

func (cr *catalogItem) Get(ctx context.Context, id string) (*entity.CatalogItem, error) {
	query := `
	SELECT id, name, price
	FROM CatalogItems
	WHERE id = ?
	LIMIT 1
	`
	row := cr.db.QueryRowContext(ctx, query, id)
	var item entity.CatalogItem
	if err := row.Scan(
		&item.Id,
		&item.Name,
		&item.Price,
	); err != nil {
		return nil, err
	}
	return &item, nil
}

func (cr *catalogItem) List(ctx context.Context) (entity.CatalogItems, error) {
	query := `
	SELECT id, name, price
	FROM CatalogItems
	`
	rows, err := cr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items entity.CatalogItems
	for rows.Next() {
		var item entity.CatalogItem
		if err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (cr *catalogItem) ListByName(ctx context.Context, name string) (entity.CatalogItems, error) {
	query := `
	SELECT id, name, price
	FROM CatalogItems
	WHERE name LIKE ?
	`
	rows, err := cr.db.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items entity.CatalogItems
	for rows.Next() {
		var item entity.CatalogItem
		if err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (cr *catalogItem) ListByIDs(ctx context.Context, ids []string) (entity.CatalogItems, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	//nolint: gosec // ignore SQL string concatenation
	query := `
	SELECT id, name, price
	FROM CatalogItems
	WHERE id IN (` + strings.Join(placeholders, ",") + `)
	`
	rows, err := cr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items entity.CatalogItems
	for rows.Next() {
		var item entity.CatalogItem
		if err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (cr *catalogItem) Create(ctx context.Context, item *entity.CatalogItem) error {
	query := `
	INSERT INTO CatalogItems (
	id, name, price
	)
	VALUES (?, ?, ?)
	`
	if _, err := cr.db.ExecContext(
		ctx,
		query,
		item.GetId(),
		item.GetName(),
		item.GetPrice(),
	); err != nil {
		return err
	}
	return nil
}

func (cr *catalogItem) Update(ctx context.Context, item *entity.CatalogItem) error {
	query := `
	UPDATE CatalogItems
	SET name = ?, price = ?
	WHERE id = ?
	`
	if _, err := cr.db.ExecContext(
		ctx,
		query,
		item.GetName(),
		item.GetPrice(),
		item.GetId(),
	); err != nil {
		return err
	}
	return nil
}

func (cr *catalogItem) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM CatalogItems
	WHERE id = ?
	`
	if _, err := cr.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
