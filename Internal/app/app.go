package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Threat struct {
	CountryCode string    `json:"countrycode"`
	ThreatLevel int       `json:"threatlevel"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pubdate"`
}

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetRoutes() {
	app.Router.
		Methods("GET").
		Path("/threats/{cc}").
		HandlerFunc(app.getThreats)
}

func (app *App) getThreats(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var threats []Threat
	result, err := app.Database.Query("SELECT COUNTRY_CODE, THREAT_LEVEL, TITLE, DESCRIPTION from RSS_THREATS")
	if err != nil {
		fmt.Println("Query failed")
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var threat Threat
		err := result.Scan(&threat.CountryCode, &threat.ThreatLevel, &threat.Title, &threat.Description)
		if err != nil {
			fmt.Println("Exits with status code 3")
			panic(err.Error())
		}
		threats = append(threats, threat)
	}
	json.NewEncoder(w).Encode(threats)
}
