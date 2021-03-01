package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bharris183/travel_alerts/innternal/app"
	"github.com/bharris183/travel_alerts/innternal/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	database, err := db.OpenConnection()
	if err != nil {
		fmt.Println("Connection open failure")
		panic(err.Error())
	}
	defer db.Close()

	app := &app.App{
		Router:   mux.NewRouter(),
		Database: database,
	}

	log.Fatal(http.ListenAndServe(":8000", app.Router))

}
