package main

import (
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
}
