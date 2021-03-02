package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bharris183/travel_alerts/internal/app"
	"github.com/bharris183/travel_alerts/internal/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	database, err := db.OpenConnection()
	if err != nil {
		fmt.Println("Connection open failure")
		panic(err.Error())
	}
	defer database.Close()

	app := &app.App{
		Router:   mux.NewRouter(),
		Database: database,
	}

	app.SetRoutes()

	log.Fatal(http.ListenAndServe(":8000", app.Router))

}
