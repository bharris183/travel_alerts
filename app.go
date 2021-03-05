package main

import (
	"database/sql"
	"fmt"
	"log"
    "net/http"
    "encoding/json"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
    
    serverName := "127.0.0.1:6033"
    dbconn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", user, password, serverName, dbname)
    var err error
    a.DB, err = sql.Open("mysql", dbconn)
    if err != nil {
        fmt.Println("Could not open database. Connection string: " + dbconn)
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()

    a.initializeRoutes()
}

func (a *App) Run(addr string) {
    log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/threat/{country}", a.getThreat).Methods("GET")
}

func (a *App) getThreat(w http.ResponseWriter, r *http.Request) {
    a.addThreat()
    vars := mux.Vars(r)
    countryCode, _ := vars["country"]

    // ToDo Verify country
    //respondWithError(w, http.StatusBadRequest, "Invalid product Country")
 
    t := threat{CountryCode: countryCode}
    if err := t.getThreat(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Threats not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, t)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (a *App) addThreat() {
	const strSql = `insert into RSS_THREATS (COUNTRY_CODE, THREAT_LEVEL, TITLE, LINK, DESCRIPTION, PUB_DATE ) 
    values ('SN', 3, 'Senegal is safe, mostly', 'http://sn.com', 'Senegal is a lovely place. Just stick to your hotel.', '2021-05-11' )`
	a.DB.Exec(strSql)
}