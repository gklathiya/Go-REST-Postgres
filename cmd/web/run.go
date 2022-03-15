package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/gklathiya/Go-REST-Postgres/internal/config"
	"github.com/gklathiya/Go-REST-Postgres/internal/driver"
	"github.com/gklathiya/Go-REST-Postgres/internal/handlers"
	"github.com/gklathiya/Go-REST-Postgres/internal/helpers"
	"github.com/gklathiya/Go-REST-Postgres/internal/render"
)

const portNumber = ":2450"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func run() (*driver.DB, error) {

	dbName := flag.String("dbname", "", "Database Name")
	dbHost := flag.String("dbhost", "localhost", "Database Host")
	dbUser := flag.String("dbuser", "", "Database User")
	dbPass := flag.String("dbpass", "", "Database Password")
	jwtSec := flag.String("jwtsec", "", "JWT Secret Key")
	dbPort := flag.String("dbport", "5432", "Database Port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL Setting(disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	app.JWTKey = *jwtSec
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	app.InProduction = false

	// connect to database
	log.Println("Connecting to database....")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	//"host=localhost port=5432 dbname=bookings user=postgres password=root"
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Connecting DB Failed", err)
	}
	log.Println("Connected to the Database !")

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}
