package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/entity"
	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/customer/repository"
)

type customerRepository struct {
	db SQLExecutor
}

func NewCustomerRepository(db *sql.DB) repository.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (cr *customerRepository) Get(ctx context.Context, id string) (*entity.Customer, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, name, email, street, city, country
	FROM Customers
	WHERE id = ?
	LIMIT 1
	`

	row := executor.QueryRowContext(ctx, query, id)
	var customer entity.Customer
	if err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Street,
		&customer.City,
		&customer.Country,
	); err != nil {
		return nil, err
	}
	return &customer, nil
}

func (cr *customerRepository) List(ctx context.Context) ([]entity.Customer, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	SELECT id, name, email, street, city, country
	FROM Customers
	`

	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []entity.Customer
	for rows.Next() {
		var customer entity.Customer
		if err = rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Street,
			&customer.City,
			&customer.Country,
		); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (cr *customerRepository) Create(ctx context.Context, customer entity.Customer) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	INSERT INTO Customers (
	id, name, email, street, city, country
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Street,
		customer.City,
		customer.Country,
	); err != nil {
		return err
	}
	return nil
}

func (cr *customerRepository) Update(ctx context.Context, customer entity.Customer) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	UPDATE Customers
	SET name = ?, email = ?, street = ?, city = ?, country = ?
	WHERE id = ?
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		customer.Name,
		customer.Email,
		customer.Street,
		customer.City,
		customer.Country,
		customer.ID,
	); err != nil {
		return err
	}
	return nil
}

func (cr *customerRepository) Delete(ctx context.Context, id string) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `
	DELETE FROM Customers
	WHERE id = ?
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
