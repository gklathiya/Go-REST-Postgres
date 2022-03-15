package main

import (
	"fmt"
	"log"
	"net/http"
)

// Main is the main application function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	log.Println(fmt.Sprintf("Application Started on: http://localhost%s", portNumber))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Server faild to Start", err)
	}
}
