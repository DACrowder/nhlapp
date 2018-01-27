package main

import (
	"log"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Db - database pointer to main storage database
var Db *sqlx.DB
var connStr = "user=Doyle dbname=nhlapp sslmode=disable"

// IsUniqueViolation returns true if the supplied error resulted from
// a unique constraint violation
// thanks for the function Nick
func IsUniqueViolation(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		return err.Code == "23505"
	}

	return false
}

func main() {
	var err error
	Db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		//cannot connect to database
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/shiftapi/v1/{game_id}", displayGame).Methods("GET")

	//http.ListenAndServe(":9999", r)

	scrape("2017020028")

	GetEvents("2017020028")

	Db.Close()
}
