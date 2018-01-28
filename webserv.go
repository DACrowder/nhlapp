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

	category := "SHOT"

	gameID, ok := vars["game_id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playerID, ok := vars["player_id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	scrape(gameID)
	GetEvents(gameID)

	q := `SELECT * FROM event WHERE game_id = ($1) AND
			player_id = ($2) AND event_type = ($3)`

	rows, err := Db.Query(q, gameID, playerID, category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for rows.Next() {
		eventOut := &Event{}
		rows.Scan(&eventOut)
		fmt.Fprintf(w, "%v#\n", eventOut)
	}

}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	gameID, ok := vars["game_id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	scrape(gameID)
	GetEvents(gameID)

	type players struct {
		PlayerID []int `json:"playerId"`
	}

	q := `SELECT DISTINCT player_id FROM shift WHERE game_id = $1`

	rows, err := Db.Query(q, gameID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for rows.Next() {
		output := players{}
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
