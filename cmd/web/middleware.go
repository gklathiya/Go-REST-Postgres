package main

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gklathiya/Go-REST-Postgres/internal/helpers"
	"github.com/justinas/nosurf"
)

type Data struct {
	Message string `json:"message"`
	JWT     string `json:"jwt"`
}

//NoSurf adds CSRF protection to POST evry request
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler

}

//IsAuthorized check whether user is authorized or not
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") == "" {
			resp := Data{
				Message: "JWT Token is not set into request header",
				JWT:     r.Header.Get("authorization"),
			}
			out, err := json.MarshalIndent(resp, "", "\t")
			if err != nil {
				helpers.ServerError(w, err) // using our own error handler
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(out)
			return
		}

		var mySigningKey = []byte(app.JWTKey)

		token, err := jwt.Parse(r.Header.Get("authorization"), func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, nil
			}
			return mySigningKey, nil
		})
		resp := Data{
			Message: "Wrong JWT or Expired Token",
			JWT:     r.Header.Get("authorization"),
		}
		if err != nil {
			out, err := json.MarshalIndent(resp, "", "\t")
			if err != nil {
				helpers.ServerError(w, err) // using our own error handler
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(out)
			return
		}
		if token.Valid {
			handler.ServeHTTP(w, r)
			return
		}
	}
}
