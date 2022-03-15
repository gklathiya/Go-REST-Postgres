package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gklathiya/Go-REST-Postgres/internal/config"
	"github.com/gklathiya/Go-REST-Postgres/internal/driver"
	"github.com/gklathiya/Go-REST-Postgres/internal/helpers"
	"github.com/gklathiya/Go-REST-Postgres/internal/models"
	"github.com/gklathiya/Go-REST-Postgres/internal/repository"
	"github.com/gklathiya/Go-REST-Postgres/internal/repository/dbrepo"
	"github.com/go-chi/chi"
)

var Repo *Repository

type Data struct {
	Message string `json:"message"`
}
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// getProduct is the return product in json from database
func (m *Repository) GetProduct(w http.ResponseWriter, r *http.Request) {

	products, err := m.DB.GetProducts()

	if err != nil {
		helpers.ServerError(w, err) // using our own error handler
		return
	}
	//log.Println(products)

	out, err := json.MarshalIndent(products, "", "\t")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) AddProduct(w http.ResponseWriter, r *http.Request) {

	var prod models.Product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.InsertProduct(prod)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := Data{
		Message: "Product Inserted Successfuly",
	}
	out, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var prod models.Product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	prod_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	result, err := m.DB.CheckProductAvailability(prod_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	resp := Data{}
	if result {
		err = m.DB.UpdateProduct(prod_id, prod)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		resp = Data{
			Message: "Product Updated Successfuly",
		}
	} else {
		resp = Data{
			Message: "Product not found for the given ID",
		}
	}

	out, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		helpers.ServerError(w, err) // using our own error handler
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	prod_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	result, err := m.DB.CheckProductAvailability(prod_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	resp := Data{}
	if result {
		res, err := m.DB.DeleteProduct(prod_id)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		if res {
			resp = Data{
				Message: "Product Deleted Successfuly",
			}
		}
	} else {
		resp = Data{
			Message: "Product not found for the given ID",
		}
	}

	out, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		helpers.ServerError(w, err) // using our own error handler
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) SignIn(w http.ResponseWriter, r *http.Request) {

	var authDetails Authentication

	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		helpers.ServerError(w, err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	res, err := m.DB.CheckUserAutherisation(authDetails.Email, authDetails.Password)

	if !res {
		helpers.ServerError(w, err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := helpers.GenerateJWT(authDetails.Email)
	if err != nil {
		helpers.ServerError(w, err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var token Token
	token.Email = authDetails.Email
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(token)
	out, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Write(out)
}
