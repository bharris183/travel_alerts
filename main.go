package main
import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Threat struct {
	CountryCode string `json:"countrycode"`
	ThreatLevel int `json:"threatlevel"`
	Title string `json:"title"`
	Link string `json:"link"`
	Description string `json:"description"`
	PubDate time.Time `json:"pubdate"`
  }

var db *sql.DB
var err error

func main() {
  db, err = sql.Open("mysql", "root:sw0rdfish@tcp(127.0.0.1:6033)/threat_alerts")
  if err != nil {
    fmt.Println("Exits with status code 1")
	panic(err.Error())
  }
  defer db.Close()

  router := mux.NewRouter()

  router.HandleFunc("/threats/{cc}", getThreats).Methods("GET")

  http.ListenAndServe(":8000", router)
/*
  _, err := db.Query("CREATE DATABASE IF NOT EXISTS threat_alerts CHARACTER SET latin1 COLLATE latin1_swedish_ci")
  if err != nil {
	panic(err.Error())
  }
  */
}

func getThreats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var threats []Threat
	result, err := db.Query("SELECT COUNTRY_CODE, THREAT_LEVEL, TITLE, DESCRIPTION from RSS_THREATS")
	if err != nil {
		fmt.Println("Exits with status code 2")
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