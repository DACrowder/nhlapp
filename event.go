package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// thanks to https://mholt.github.io/json-to-go/ for helping create this struct
type Event struct {
	GamePk   int    `json:"gamePk"` // game_id
	Link     string `json:"link"`   // game link
	GameData struct {
		Game struct {
			Pk int `json:"pk"` // game_id
		} `json:"game"`
	} `json:"gameData"`
	LiveData struct {
		Plays struct {
			AllPlays []struct {
				Result struct {
					Event       string `json:"event"`
					EventCode   string `json:"eventCode"`
					EventTypeID string `json:"eventTypeId"` //event type
					Description string `json:"description"`
				} `json:"result"`
				About struct {
					EventIdx            int       `json:"eventIdx"` // event_id
					EventID             int       `json:"eventId"`
					Period              int       `json:"period"` // period
					PeriodType          string    `json:"periodType"`
					OrdinalNum          string    `json:"ordinalNum"`
					PeriodTime          string    `json:"periodTime"` // period_time
					PeriodTimeRemaining string    `json:"periodTimeRemaining"`
					DateTime            time.Time `json:"dateTime"`
					Goals               struct {
						Away int `json:"away"`
						Home int `json:"home"`
					} `json:"goals"`
				} `json:"about"`
				Coordinates struct { //coords need to test
					X float32 `json:"x"` // coord_x originally not in generated struct
					Y float32 `json:"y"` //coord_y originally not in generated struct
				} `json:"coordinates"`
				Players []struct { // player 1 = Players[0], player 2 = Players[-1]
					Player struct {
						ID       int    `json:"id"` // playerx_id
						FullName string `json:"fullName"`
						Link     string `json:"link"`
					} `json:"player"`
					PlayerType string `json:"playerType"` // playerx_type
				} `json:"players,omitempty"`
			} `json:"allPlays"`
		} `json:"plays"`
	} `json:"liveData"`
}

func GetEvents(gameID string) (*Event, error) {
	apiURL := fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/game/%s/feed/live", gameID)

	client := &http.Client{}
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	data := Event{}
	dataDec := json.NewDecoder(response.Body)
	err = dataDec.Decode(&data)
	if err != nil {
		return nil, err
	}

	for _, cur := range data.LiveData.Plays.AllPlays {
		fmt.Printf("idx: %d, code: %s, x: %f, y: %f\n", cur.About.EventIdx, cur.Result.EventTypeID, cur.Coordinates.X, cur.Coordinates.Y)
	}

	return &data, nil
}
