package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/microservice-k8s-demo/catalog/entity"
	"github.com/tusmasoma/microservice-k8s-demo/catalog/repository"
)

type catalogItemRepository struct {
	db SQLExecutor
}

func NewCatalogItemRepository(db *sql.DB) repository.CatalogItemRepository {
	return &catalogItemRepository{
		db: db,
	}
}

func (cr *catalogItemRepository) Get(ctx context.Context, id string) (*entity.CatalogItem, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, name, price
	FROM CatalogItems
	WHERE id = ?
	LIMIT 1
	`

	row := executor.QueryRowContext(ctx, query, id)
	var item entity.CatalogItem
	if err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Price,
	); err != nil {
		return nil, err
	}
	return &item, nil
}

func (cr *catalogItemRepository) List(ctx context.Context) ([]entity.CatalogItem, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, name, price
	FROM CatalogItems
	`

	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.CatalogItem
	for rows.Next() {
		var item entity.CatalogItem
		if err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (cr *catalogItemRepository) ListByName(ctx context.Context, name string) ([]entity.CatalogItem, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, name, price
	FROM CatalogItems
	WHERE name = ?
	`

	rows, err := executor.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.CatalogItem
	for rows.Next() {
		var item entity.CatalogItem
		if err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (cr *catalogItemRepository) ListByNameContaining(ctx context.Context, name string) ([]entity.CatalogItem, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, name, price
	FROM CatalogItems
	WHERE name LIKE ?
	`

	rows, err := executor.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.CatalogItem
	for rows.Next() {
		var item entity.CatalogItem
		if err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (cr *catalogItemRepository) Create(ctx context.Context, item entity.CatalogItem) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	INSERT INTO CatalogItems (
	id, name, price
	)
	VALUES (?, ?, ?)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		item.ID,
		item.Name,
		item.Price,
	); err != nil {
		return err
	}
	return nil
}

func (cr *catalogItemRepository) Update(ctx context.Context, item entity.CatalogItem) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	UPDATE CatalogItems
	SET name = ?, price = ?
	WHERE id = ?
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		item.Name,
		item.Price,
		item.ID,
	); err != nil {
		return err
	}
	return nil
}

func (cr *catalogItemRepository) Delete(ctx context.Context, id string) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	DELETE FROM CatalogItems
	WHERE id = ?
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
