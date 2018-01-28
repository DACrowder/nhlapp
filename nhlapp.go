package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	ConnStr string
}

// Db - database pointer to main storage database
var Db *sqlx.DB

// IsUniqueViolation returns true if the supplied error resulted from
// a unique constraint violation
// thanks for the function Nick
func IsUniqueViolation(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		return err.Code == "23505"
	}

	return false
}

// TimeConvert - converts time strings from nhl API to time Duration objects
func TimeConvert(timeString string) (int, error) {
	clock := strings.Split(timeString, ":")

	timeFull := fmt.Sprintf("%s%s", clock[0], clock[1])

	return strconv.Atoi(timeFull)
}

func main() {
	conf := Configuration{}
	err := gonfig.GetConf("config.json", &conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", conf)

	Db, err = sqlx.Connect("postgres", conf.ConnStr)
	if err != nil {
		//cannot connect to database
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/shiftapi/v1/{game_id}", getPlayers).Methods("GET")
	r.HandleFunc("/shiftapi/v1/{game_id}/player/{player_id}", displayGame).Methods("GET")

	http.ListenAndServe(":9999", r)

	//scrape("2017020028")

	//GetEvents("2017020028")

	Db.Close()
}
