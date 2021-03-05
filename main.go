package main

import (
	//"os"
)

func main() {
	a := App{}
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
	a.Run(":8010")



}