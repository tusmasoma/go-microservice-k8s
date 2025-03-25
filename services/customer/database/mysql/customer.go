package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-microservice-k8s/services/customer/database"
	"github.com/tusmasoma/go-microservice-k8s/services/customer/entity"
)

type customer struct {
	db *sql.DB
}

func NewCustomer(db *sql.DB) database.Customer {
	return &customer{
		db: db,
	}
}

func (cr *customer) Get(ctx context.Context, id string) (*entity.Customer, error) {
	query := `
	SELECT id, name, email, street, city, country
	FROM Customers
	WHERE id = ?
	LIMIT 1
	`
	row := cr.db.QueryRowContext(ctx, query, id)
	var customer entity.Customer
	if err := row.Scan(
		&customer.Id,
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

func (cr *customer) List(ctx context.Context) (entity.Customers, error) {
	query := `
	SELECT id, name, email, street, city, country
	FROM Customers
	`
	rows, err := cr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers entity.Customers
	for rows.Next() {
		var customer entity.Customer
		if err = rows.Scan(
			&customer.Id,
			&customer.Name,
			&customer.Email,
			&customer.Street,
			&customer.City,
			&customer.Country,
		); err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}

func (cr *customer) Create(ctx context.Context, customer *entity.Customer) error {
	query := `
	INSERT INTO Customers (
	id, name, email, street, city, country
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	if _, err := cr.db.ExecContext(
		ctx,
		query,
		customer.GetId(),
		customer.GetName(),
		customer.GetEmail(),
		customer.GetStreet(),
		customer.GetCity(),
		customer.GetCountry(),
	); err != nil {
		return err
	}
	return nil
}

func (cr *customer) Update(ctx context.Context, customer *entity.Customer) error {
	query := `
	UPDATE Customers
	SET name = ?, email = ?, street = ?, city = ?, country = ?
	WHERE id = ?
	`
	if _, err := cr.db.ExecContext(
		ctx,
		query,
		customer.GetName(),
		customer.GetEmail(),
		customer.GetStreet(),
		customer.GetCity(),
		customer.GetCountry(),
		customer.GetId(),
	); err != nil {
		return err
	}
	return nil
}

func (cr *customer) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM Customers
	WHERE id = ?
	`
	if _, err := cr.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
