package main

import (
	"fmt"
	"log"
)

// EventRoster - record to match each player on the ice to an event
type EventRoster struct {
	GameID   int    `json:"gameId" db:"game_id"`
	EventID  int    `json:"eventId" db:"event_id"`
	Team     string `json:"team" db:"team"`
	PlayerID int
}

// CreateEventRoster - links all players on the ice to each event
func CreateEventRoster(gameID string) {
	//LOL THIS QUERY WORKED
	q := `INSERT INTO event_roster (game_id, event_id, team, player_id) 
			SELECT e.game_id, e.event_id, s.team, s.player_id
			FROM event AS e, shift AS s
			WHERE e.game_id = ($1) AND
			s.period = e.period AND
			s.time_start <= e.period_time AND
			s.time_end > e.period_time`

	result, err := Db.Exec(q, gameID)
	if err != nil {
		if !IsUniqueViolation(err) {
			log.Println(err)
			return
		}
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%d row(s) created.\n", count)

}
