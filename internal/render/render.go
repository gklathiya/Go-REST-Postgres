package render

import "github.com/gklathiya/Go-REST-Postgres/internal/config"

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}
