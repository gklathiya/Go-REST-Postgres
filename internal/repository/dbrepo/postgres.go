package dbrepo

import (
	"context"
	"time"

	"github.com/gklathiya/Go-REST-Postgres/internal/models"
)

func (m *postgresDBRepo) GetProducts() ([]models.Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // This is will check is transaction not done within 3 seconds then cancel it
	defer cancel()

	var products []models.Product
	query := `SELECT *
			FROM 
			products`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return products, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Product
		err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
		)
		if err != nil {
			return products, err
		}
		products = append(products, i)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil
}

func (m *postgresDBRepo) InsertProduct(prod models.Product) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // This is will check is transaction not done within 3 seconds then cancel it
	defer cancel()

	stmt := `INSERT INTO products(name, price)
			VALUES ($1,$2)`

	_, err := m.DB.QueryContext(ctx, stmt, prod.Name, prod.Price)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) CheckProductAvailability(prod_id int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // This is will check is transaction not done within 3 seconds then cancel it
	defer cancel()

	stmt := `SELECT COUNT(id) FROM products WHERE id=$1`
	var numRows int
	err := m.DB.QueryRowContext(ctx, stmt,
		prod_id,
	).Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return false, nil
	}
	return true, nil
}

func (m *postgresDBRepo) UpdateProduct(prod_id int, prod models.Product) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // This is will check is transaction not done within 3 seconds then cancel it
	defer cancel()

	query := `UPDATE products SET name=$1, price=$2
			WHERE
				id = $3`

	_, err := m.DB.ExecContext(ctx, query,
		prod.Name,
		prod.Price,
		prod_id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteProduct(prod_id int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // This is will check is transaction not done within 3 seconds then cancel it
	defer cancel()

	query := `DELETE FROM products WHERE
				id = $1`
	_, err := m.DB.ExecContext(ctx, query, prod_id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *postgresDBRepo) CheckUserAutherisation(email, password string) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT COUNT(id) FROM users WHERE email=$1 AND password=$2`
	var numRows int
	err := m.DB.QueryRowContext(ctx, stmt,
		email,
		password,
	).Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return false, nil
	}
	return true, nil
}
