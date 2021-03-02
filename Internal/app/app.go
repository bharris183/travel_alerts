package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"github.com/ungerik/go-rss"

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
	Threats []Threat
}

func (app *App) SetRoutes() {
	app.Router.
		Methods("GET").
		Path("/threats/{cc}").
		HandlerFunc(app.getThreats)

	app.Router.
		Methods("GET").
		Path("/rssthreats").
		HandlerFunc(app.getThreatsFromRss)
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

func (app *App) getThreatsFromRss(w http.ResponseWriter, r *http.Request) {

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

	var threats []Threat
	// Now let's make our Threats with only the data we need
	for _, item := range channel.Item {
		var threat Threat
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
	for _, t := range threats {
		fmt.Fprintf(w, strconv.Itoa(t.ThreatLevel) + "\n")
	}
	//return threats
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