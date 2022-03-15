package main

import (
	"net/http"

	"github.com/gklathiya/Go-REST-Postgres/internal/config"
	"github.com/gklathiya/Go-REST-Postgres/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)

	mux.Get("/getProducts", IsAuthorized(handlers.Repo.GetProduct))
	mux.Post("/addProduct", handlers.Repo.AddProduct)
	mux.Put("/editProduct/{id}", handlers.Repo.UpdateProduct)
	mux.Delete("/Product/{id}", handlers.Repo.DeleteProduct)
	mux.Post("/sigin", handlers.Repo.SignIn)
	return mux
}
