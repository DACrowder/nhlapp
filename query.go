package main

import (
	"fmt"
)

func getTeamShots(gameID string, team string) (int, error) {
	if gameID == "" || team == "" {
		return 0, fmt.Errorf("game or team not given")
	}
	q := `select count (*) from event where game_id = $1 and player1_team = $2 and event_type = 'SHOT' or event_type = 'GOAL'`

	row := Db.QueryRow(q, gameID, team)
	/*if err != nil {
		return 0, err
	}*/

	var count int
	row.Scan(&count)

	return count, nil
}
