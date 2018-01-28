package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Player struct {
	Copyright string `json:"copyright"`
	People    []struct {
		ID               int    `json:"id"`
		FullName         string `json:"fullName"`
		Link             string `json:"link"`
		FirstName        string `json:"firstName"`
		LastName         string `json:"lastName"`
		PrimaryNumber    string `json:"primaryNumber"`
		BirthDate        string `json:"birthDate"`
		CurrentAge       int    `json:"currentAge"`
		BirthCity        string `json:"birthCity"`
		BirthCountry     string `json:"birthCountry"`
		Nationality      string `json:"nationality"`
		Height           string `json:"height"`
		Weight           int    `json:"weight"`
		Active           bool   `json:"active"`
		AlternateCaptain bool   `json:"alternateCaptain"`
		Captain          bool   `json:"captain"`
		Rookie           bool   `json:"rookie"`
		ShootsCatches    string `json:"shootsCatches"`
		RosterStatus     string `json:"rosterStatus"`
		CurrentTeam      struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"currentTeam"`
		PrimaryPosition struct {
			Code         string `json:"code"`
			Name         string `json:"name"`
			Type         string `json:"type"`
			Abbreviation string `json:"abbreviation"`
		} `json:"primaryPosition"`
	} `json:"people"`
}

func getPlayerPosition(playerID int) (string, error) {
	apiURL := fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/people/%s", playerID)

	client := &http.Client{}
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	data := Player{}
	dataDec := json.NewDecoder(response.Body)
	err = dataDec.Decode(&data)
	if err != nil {
		return "", err
	}

	return data.People[0].PrimaryPosition.Abbreviation, nil
}
