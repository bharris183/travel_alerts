package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"encoding/json"
	"reflect"
	"testing"   
    "log"
)

var a App

func TestMain(m *testing.M) {
    a.Initialize(
		/*
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
		*/
		"sa_threatalerts",
		"sw0rdfish",
		"threat_alerts",
	)

    ensureTableExists()
    code := m.Run()
    clearTable()
    os.Exit(code)
}


func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM RSS_THREATS")
}

const tableCreationQuery = `CREATE  TABLE IF NOT EXISTS threat_alerts.RSS_THREATS 
(
	COUNTRY_CODE CHAR(2) NOT NULL,
	THREAT_LEVEL INT NOT NULL ,
	TITLE VARCHAR(255) ,
	LINK VARCHAR(255) ,
	DESCRIPTION VARCHAR(255) , 
	PUB_DATE DATE ,
	PRIMARY KEY (COUNTRY_CODE)
)`

func TestGetThreatForNonesxistentCountry(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/threat/json/XX", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusBadRequest, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Threat not found for country XX" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Threat not found for country XX' . Got '%s'", m["error"])
    }
}

func TestGetThreat(t *testing.T) {
    clearTable()
    addThreat()

    req, _ := http.NewRequest("GET", "/threat/json/SN", nil)
    response := executeRequest(req)

	fmt.Println(response)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "" {
        t.Errorf("Expected no error. Got '%s'", m["error"])
	}
	checkExpectedResponse(t, m)
}

func checkExpectedResponse(t *testing.T, m map[string]string) {
	var want, got string
	var wantGot = "Expected '%s', got %s"
	want = "SN"
	got = m["countrycode"] 
	if got != want {
		t.Errorf(wantGot, want, got)
	}
	//want = "3"
	got = m["threatlevel"]
	fmt.Println(reflect.TypeOf(got).String())
	//if got != want {
	//	t.Errorf(wantGot, want, got)
	//}
	want = "Senegal is safe, mostly"
	got = m["title"] 
	if got != want {
		t.Errorf(wantGot, want, got)
	}
	want = "Senegal is a lovely place. Just stick to your hotel."
	got = m["description"] 
	if got != want {
		t.Errorf(wantGot, want, got)
	}
	want = "2021-05-11T00:00:00Z"
	got = m["pubdate"] 
	if got != want {
		t.Errorf(wantGot, want, got)
	}
}

func addThreat() {
	const strSql = `insert into RSS_THREATS (COUNTRY_CODE, THREAT_LEVEL, TITLE, LINK, DESCRIPTION, PUB_DATE ) 
    values ('SN', 3, 'Senegal is safe, mostly', 'http://sn.com', 'Senegal is a lovely place. Just stick to your hotel.', '2021-05-11' )`
	a.DB.Exec(strSql)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}