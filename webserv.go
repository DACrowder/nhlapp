package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func displayGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	gameID := vars["game_id"]
	if gameID == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	scrape(gameID)
	GetEvents(gameID)

	type player struct {
		PlayerID int `json:"playerId"`
	}

	q := `SELECT DISTINCT player_id FROM shift WHERE game_id = $1`

	rows, err := Db.Query(q, gameID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for rows.Next() {
		output := player{}
		err = rows.Scan(&output.PlayerID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonOut, err := json.Marshal(output)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%s\n", string(jsonOut))
	}

}
