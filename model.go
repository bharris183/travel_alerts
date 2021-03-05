package main

import (
    "database/sql"
    "errors"
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

const sqlGetProduct = `select COUNTRY_CODE, THREAT_LEVEL, TITLE, 
	LINK, DESCRIPTION, PUB_DATE from RSS_THREATS 
	where COUNTRY_CODE = ?` 


func (t *threat) getThreat(db *sql.DB) error {
    return db.QueryRow(sqlGetProduct, t.CountryCode).Scan(
		&t.CountryCode, &t.ThreatLevel, &t.Title, &t.Link, &t.Description, &t.PubDate)
}


func (t *threats) loadThreats(db *sql.DB) error {
	return errors.New("Not implemented")
  }


