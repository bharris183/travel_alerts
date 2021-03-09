package main

import (
	"database/sql"
	"fmt"
	"log"
    "net/http"
    "encoding/json"
    "strconv"
    "strings"
    "time"

	"github.com/gorilla/mux"
    "github.com/ungerik/go-rss"
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
    //a.Router.HandleFunc("/rss", a.getThreatsFromRss).Methods("GET")
    a.Router.HandleFunc("/update", a.updateThreatsFromRss).Methods("GET")
}

// Get a single threat from a country code
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

// Loading our table

func (a *App) getThreatsFromRss() ( []threat ) {

	rssUrl := "https://travel.state.gov/_res/rss/TAsTWs.xml"

	// Similar to reading a file, only this is rss
	resp, err := rss.Read(rssUrl, false) // In rss package true for reddit feeds
	if err != nil {
		fmt.Println(err)
	}

	// Result is a Channel struct representing the rss data
	channel, err := rss.Regular(resp)
	if err != nil {
		fmt.Println(err)
	}

	var threats []threat
	// Now let's make our Threats with only the data we need
	for _, item := range channel.Item {
		var threat threat
		title := item.Title
		link := item.Link
		description := item.Description
		var pubDate time.Time
		pd, err := item.PubDate.ParseWithFormat("Mon, 02 Jan 2006")
		if err != nil {
			pubDate = time.Now()
		} else {
			pubDate = pd
		}
		countryCode := getCountryCode(item.Category)
		threatLevel := getThreatLevel(item.Category)
		
		threat.Title = title
		threat.Link = link
		threat.Description = description
		threat.PubDate = pubDate
		threat.CountryCode = countryCode
		if i, err := strconv.Atoi(threatLevel); err == nil {
			threat.ThreatLevel = i
		} else {
			threat.ThreatLevel = 0
		}
		threats = append(threats, threat)
	}
	return threats
}

func getThreatLevel(category []string) string {
	for _, s := range category {
		if strings.HasPrefix(s, "Level") {
			strParts := strings.Split(s, ":")
			s1 := strParts[0]
			strParts = strings.Split(s1, " ")
			return strParts[1]
		}
	}
	return ""
}

func getCountryCode(category []string) string {
	for _, s := range category {
		if len(s) == 2 {
			return s
		}
	}
	return ""
}

func (a *App) updateThreatsFromRss(w http.ResponseWriter, r *http.Request) {
    t := threats{LastUpdated: time.Now()}
    t.Threats = a.getThreatsFromRss()
    rowsUpdated, err := t.loadThreatsInDatbase(a.DB)
    if err != nil {
        fmt.Printf("Error loading threats from rss: %s", err)
    } else {
        fmt.Printf("%s rows updated", strconv.Itoa(rowsUpdated))
    }  
}