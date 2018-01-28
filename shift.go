package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type shift struct {
	Data []struct {
		DetailCode       int         `json:"detailCode"`
		Duration         string      `json:"duration"`
		EndTime          string      `json:"endTime" db:"time_end"`
		EventDescription interface{} `json:"eventDescription"`
		EventDetails     interface{} `json:"eventDetails"`
		EventNumber      int         `json:"eventNumber"`
		FirstName        string      `json:"firstName"`
		GameID           int         `json:"gameId"`
		HexValue         string      `json:"hexValue"`
		LastName         string      `json:"lastName"`
		Period           int         `json:"period" db:"period"`
		PlayerID         int         `json:"playerId" db:"player_start"`
		ShiftNumber      int         `json:"shiftNumber"`
		StartTime        string      `json:"startTime" db:"time_start"`
		TeamAbbrev       string      `json:"teamAbbrev"`
		TeamID           int         `json:"teamId"`
		TeamName         string      `json:"teamName"`
		TypeCode         int         `json:"typeCode"`
	} `json:"data"`
	Total int `json:"total"`
}

// Scrape pull shift data from nhl API
func scrape(gameID string) {
	apiURL := fmt.Sprintf("http://www.nhl.com/stats/rest/shiftcharts?cayenneExp=gameId=%s", gameID)
	fmt.Println(apiURL)
	client := &http.Client{}

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}

	data := shift{}
	dataDec := json.NewDecoder(response.Body)
	err = dataDec.Decode(&data)
	if err != nil {
		log.Println(err)
		return
	}

	for _, shiftSlice := range data.Data {
		q := `INSERT INTO shift (game_id, player_id, period, time_start, time_end)
						VALUES ($1, $2, $3, $4, $5)`
		startInt, err := TimeConvert(shiftSlice.StartTime)
		if err != nil {
			log.Println(err)
			return
		}
		stopInt, err := TimeConvert(shiftSlice.EndTime)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = Db.Exec(q, shiftSlice.GameID, shiftSlice.PlayerID, shiftSlice.Period, startInt, stopInt)
		if err != nil {
			if IsUniqueViolation(err) {
				continue
			}
			log.Fatal(err)
			return
		}
	}

}
