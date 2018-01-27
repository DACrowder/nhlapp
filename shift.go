package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type shift struct {
	Data []struct {
		DetailCode       int         `json:"detailCode"`
		Duration         string      `json:"duration"`
		EndTime          string      `json:"endTime"`
		EventDescription interface{} `json:"eventDescription"`
		EventDetails     interface{} `json:"eventDetails"`
		EventNumber      int         `json:"eventNumber"`
		FirstName        string      `json:"firstName"`
		GameID           int         `json:"gameId"`
		HexValue         string      `json:"hexValue"`
		LastName         string      `json:"lastName"`
		Period           int         `json:"period"`
		PlayerID         int         `json:"playerId"`
		ShiftNumber      int         `json:"shiftNumber"`
		StartTime        string      `json:"startTime"`
		TeamAbbrev       string      `json:"teamAbbrev"`
		TeamID           int         `json:"teamId"`
		TeamName         string      `json:"teamName"`
		TypeCode         int         `json:"typeCode"`
	} `json:"data"`
	Total int `json:"total"`
}

// Scrape pull shift data from nhl API
func Scrape(gameID string) {
	apiURL := fmt.Sprintf("http://www.nhl.com/stats/rest/shiftcharts?cayenneExp=gameId=%s", gameID)
	fmt.Println(apiURL)
	client := &http.Client{}

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	data := shift{}
	dataDec := json.NewDecoder(response.Body)
	dataDec.Decode(&data)

	for _, data := range data.Data {
		fmt.Printf("p: %d, s: %s, d: %s\n", data.PlayerID, data.StartTime, data.EndTime)
		q := `INSERT INTO shift (player_id, period, time_start, time_end)	
					VALUES ($1, $2, $3, $4)`
		result, err := Db.Exec(q, data.PlayerID, data.Period, data.StartTime, data.EndTime)
		if err != nil {
			log.Println(err)
			return
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Number of shifts recorded = %d\n", count)
	}

}
