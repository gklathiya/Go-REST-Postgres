package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gklathiya/Go-REST-Postgres/internal/config"
	"golang.org/x/crypto/bcrypt"
)

var app *config.AppConfig

//NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

//ClientError
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of ", status)
	http.Error(w, http.StatusText(status), status)
}

//ServerError
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// GenerateJWT Create JWT Token and return
func GenerateJWT(email string) (string, error) {
	var mySigningKey = []byte(app.JWTKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

//GeneratehashPassword take password as input and generate new hash password from it
func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash compare plain password with hash password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
