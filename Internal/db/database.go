package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func OpenConnection() (*sql.DB, error) {
	serverName := "127.0.0.1:6033"
	user := "root"
	password := "sw0rdfish"
	dbname := "threat_alerts"

	dbconn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=ut8mb4_unicode_ci&parseTime=true&multoStateents=true", user, password, serverName, dbname)
	db, err :- sql.Open("mysql", dbconn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
