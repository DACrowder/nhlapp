package main

import (
	"log"
)

type line struct {
	GameID      int    `db:"game_id"`
	LinePlayers string `db:"line_players"`
}

func buildLines(gameID string) error {
	q := `select distinct player1_team from event where game_id = $1 and player1_team != '';`
	//get both teams
	rows, err := Db.Queryx(q, gameID)
	if err != nil {
		return err
	}

	lines := Lines{}

	if rows.Next() {
		if err = rows.Scan(&lines.Team1Name); err != nil {
			return err
		}
	} else {
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&lines.Team2Name); err != nil {
			return err
		}
	} else {
		return err
	}

	q = `select count(*) as data_count, players as line_array
    from (select array_agg(q.player_id) as players, event_id
        from (select distinct roster.player_id, roster.event_id
                from event_roster as roster,
                    (select distinct event_id, player1_team as team
                        from event
                        where game_id = $1 and player1_team = $2
                    ) as sub_q
                where roster.game_id = $1 and roster.event_id = sub_q.event_id and roster.team = sub_q.team
                order by roster.event_id asc) as q
        group by q.event_id) as q2
    group by q2.players order by count(*)`

	/* get first team */
	rows, err = Db.Queryx(q, gameID, lines.Team1Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.Team1Line = append(lines.Team1Line, line)
	}

	/* get second game */
	rows, err = Db.Queryx(q, gameID, lines.Team2Name)

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.Team2Line = append(lines.Team2Line, line)
	}

	for _, lineSlice := range lines.Team1Line {
		q := `INSERT INTO line (game_id, line_players, team)
				VALUES ($1, $2, $3)`
		_, err := Db.Exec(q, gameID, lineSlice.LineTmp, lines.Team1Name)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	for _, lineSlice := range lines.Team2Line {
		q := `INSERT INTO line (game_id, line_players, team)
				VALUES ($1, $2, $3)`
		_, err := Db.Exec(q, gameID, lineSlice.LineTmp, lines.Team2Name)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
