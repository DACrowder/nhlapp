package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Db - database pointer to main storage database
var Db *sqlx.DB
var connStr = "user=Doyle dbname=nhlapp sslmode=disable"

func main() {
	var err error
	Db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		//cannot connect to database
		log.Fatal(err)
	}
	scrape("2017020028")

	Db.Close()

}
