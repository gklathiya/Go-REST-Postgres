package repository

import "github.com/gklathiya/Go-REST-Postgres/internal/models"

type DatabaseRepo interface {
	GetProducts() ([]models.Product, error)
	InsertProduct(prod models.Product) error
	CheckProductAvailability(prod_id int) (bool, error)
	UpdateProduct(prod_id int, prod models.Product) error
	DeleteProduct(prod_id int) (bool, error)
	CheckUserAutherisation(email, password string) (bool, error)
}
