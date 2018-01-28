package main

import (
	"fmt"
)

func getTeamShots(gameID string, team string) (int, error) {
	if gameID == "" || team == "" {
		return 0, fmt.Errorf("game or team not given")
	}
	q := `select count (*) from event where game_id = $1 and player1_team = $2 and (event_type = 'SHOT' or event_type = 'GOAL')`

	row := Db.QueryRow(q, gameID, team)
	/*if err != nil {
		return 0, err
	}*/

	var count int
	row.Scan(&count)

	return count, nil
}

type LineData struct {
	DataCount int    `db:"data_count"`
	LineTmp   string `db:"line_array"`
	Line      []int
}

func getLineShots(gameID string) error {
	q := `select count(*) as data_count, players as line_array
    from (select array_agg(q.player_id) as players, event_id
        from (select distinct roster.player_id, roster.event_id
                from event_roster as roster,
                    (select distinct event_id, player1_team as team
                        from event
                        where game_id = $1 and (event_type = 'SHOT' or event_type = 'GOAL') and player1_team = $2
                    ) as sub_q
                where roster.game_id = $1 and roster.event_id = sub_q.event_id and roster.team = sub_q.team
                order by roster.event_id asc) as q
        group by q.event_id) as q2
    group by q2.players order by count(*)`

	rows, err := Db.Queryx(q, gameID, "EDM")
	if err != nil {
		return err
	}

	fmt.Println(rows)

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		fmt.Printf("%d%v\n", line.DataCount, line.Line)
	}

	return nil
}
