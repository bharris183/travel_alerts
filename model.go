package main

import (
	"context"
    "database/sql"
	"log"
	"strings"
	"time"
)

type threats struct {
	LastUpdated time.Time     `json:"countrycode"`
	Threats		[]threat
}

type threat struct {
	CountryCode string    `json:"countrycode"`
	ThreatLevel int       `json:"threatlevel"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pubdate"`
}

const sqlGetThreat = `select COUNTRY_CODE, THREAT_LEVEL, TITLE, 
	LINK, DESCRIPTION, PUB_DATE from RSS_THREATS 
	where COUNTRY_CODE = ?` 

const sqlInsertThreat = `insert into RSS_THREATS 
	(COUNTRY_CODE, THREAT_LEVEL, TITLE, 
		LINK, DESCRIPTION, PUB_DATE) values 
	(?, ?, ?, ?, ?, ?)`

var tCache = make(map[string] threat)

func (t *threat) getThreat(db *sql.DB) error {
	tFromCache, found := tCache[t.CountryCode]
	if found {
		log.Println("Threat from cache for country ", t.CountryCode)
		*t = tFromCache
		return nil
	} else {
		err := db.QueryRow(sqlGetThreat, t.CountryCode).Scan (
			&t.CountryCode, &t.ThreatLevel, &t.Title, &t.Link, &t.Description, &t.PubDate)
		if err != nil {
			log.Println("Threat not found for country ", t.CountryCode)
			return err
		} else {
			log.Println("Threat from DB for country ", t.CountryCode)
			tCache[t.CountryCode] = *t
			return nil
		}
	}
}


func (t *threats) loadThreatsInDatbase(db *sql.DB) (rowsUpdated int, err error) {
	// For now we will not do anything with t.LastUpdated

	// Delete the cache
	tCache = make(map[string] threat)
	// Delete old entries - there are only 200 or so threats at any one time
	delete, err := db.Query("delete from RSS_THREATS")
	if err != nil {
		panic(err.Error())
	}
	defer delete.Close()

	query := "insert into RSS_THREATS (COUNTRY_CODE, THREAT_LEVEL, TITLE, LINK, DESCRIPTION, PUB_DATE) values "
	var inserts []string
	var params []interface{}
	for _, st := range t.Threats {
		inserts = append(inserts, "(?, ?, ?, ?, ?, ?)")
		params = append(params, st.CountryCode, st.ThreatLevel, st.Title, st.Link, st.Description, st.PubDate)
	}
    queryVals := strings.Join(inserts, ",")
    query = query + queryVals
    //log.Println("query is", query)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
    stmt, err := db.PrepareContext(ctx, query)
    if err != nil {
        log.Printf("Error %s when preparing SQL statement", err)
        return 0, err
    }
    defer stmt.Close()
    res, err := stmt.ExecContext(ctx, params...)
    if err != nil {
        log.Printf("Error %s when inserting row into rss_threats table", err)
        return 0, err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when finding rows affected", err)
        return 0, err
    }
    log.Printf("%d rss_threats created simulatneously", rows)
    return int(rows), nil
  }