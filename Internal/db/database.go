package db

import (
	"fmt"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func OpenConnection() (*sql.DB, error) {
	serverName := "127.0.0.1:6033"
	user := "root"
	password := "sw0rdfish"
	dbname := "threat_alerts"

	dbconn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", user, password, serverName, dbname)
	db, err := sql.Open("mysql", dbconn)
	if err != nil {
		fmt.Println("Could not open database. Connection string: " + dbconn)
		return nil, err
	}

	return db, nil
}
