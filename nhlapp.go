package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	ConnStr string
}

// Db - database pointer to main storage database
var Db *sqlx.DB
var connStr = "user=Doyle dbname=nhlapp sslmode=disable"

func main() {
	conf := Configuration{}
	err := gonfig.GetConf("config.json", &conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", conf)

	Db, err := sqlx.Connect("postgres", conf.ConnStr)
	if err != nil {
		//cannot connect to database
		log.Fatal(err)
	}

	Db.Close()

	Scrape("2017020028")

	GetEvents("2017020028")
}
