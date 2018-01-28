package main

import (
	"log"
)

type line struct {
	GameID      int    `db:"game_id"`
	LinePlayers string `db:"line_players"`
}

type winLossData struct {
	DataCount int    `db:"data_count"`
	EventID   int    `db:"event_id"`
	LineTmp   string `db:"line_array"`
}

type winLossLines struct {
	team1Name     string
	team1LineWin  []winLossData
	team2Name     string
	team2LineLoss []winLossData
	team1LineLoss []winLossData
	team2LineWin  []winLossData
}

type lineWinLoss struct {
	LinePlayers string `db:"line_players"`
	Team        string `db:"team"`
	EventID     int    `db:"event_id"`
	GameID      int    `db:"game_id"`
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
		if err = rows.Scan(&lines.team1Name); err != nil {
			return err
		}
	} else {
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&lines.team2Name); err != nil {
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
	rows, err = Db.Queryx(q, gameID, lines.team1Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.team1Line = append(lines.team1Line, line)
	}

	/* get second game */
	rows, err = Db.Queryx(q, gameID, lines.team2Name)

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.team2Line = append(lines.team2Line, line)
	}

	for _, lineSlice := range lines.team1Line {
		q := `INSERT INTO line (game_id, line_players, team)
				VALUES ($1, $2, $3)`
		_, err := Db.Exec(q, gameID, lineSlice.LineTmp, lines.team1Name)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	for _, lineSlice := range lines.team2Line {
		q := `INSERT INTO line (game_id, line_players, team)
				VALUES ($1, $2, $3)`
		_, err := Db.Exec(q, gameID, lineSlice.LineTmp, lines.team2Name)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func buildLineWinLoss(gameID string) error {
	q := `select distinct player1_team from event where game_id = $1 and player1_team != '';`
	//get both teams
	rows, err := Db.Queryx(q, gameID)
	if err != nil {
		return err
	}

	lines := winLossLines{}

	if rows.Next() {
		if err = rows.Scan(&lines.team1Name); err != nil {
			return err
		}
	} else {
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&lines.team2Name); err != nil {
			return err
		}
	} else {
		return err
	}

	q = `select count(*) as data_count, event_id, players as line_array
    from (select array_agg(q.player_id) as players, event_id
        from (select distinct roster1.player_id, roster1.event_id
                from event_roster as roster1,
                    (select distinct event_id, player1_team as team
                        from event
                        where game_id = $1 AND player1_team = $2
                    ) as sub_q
                where roster1.game_id = $1 and roster1.event_id = sub_q.event_id and roster1.team = $3
                order by roster1.event_id asc) as q
        group by q.event_id) as q2
    group by q2.players, q2.event_id order by count(*), q2.event_id`

	/* get first team */
	rows, err = Db.Queryx(q, gameID, lines.team1Name, lines.team1Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := winLossData{}
		rows.StructScan(&line)
		lines.team1LineWin = append(lines.team1LineWin, line)
	}

	/* get second game */
	rows, err = Db.Queryx(q, gameID, lines.team1Name, lines.team2Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := winLossData{}
		rows.StructScan(&line)
		lines.team2LineLoss = append(lines.team2LineLoss, line)
	}

	rows, err = Db.Queryx(q, gameID, lines.team2Name, lines.team2Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := winLossData{}
		rows.StructScan(&line)
		lines.team2LineWin = append(lines.team2LineWin, line)
	}

	rows, err = Db.Queryx(q, gameID, lines.team2Name, lines.team1Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := winLossData{}
		rows.StructScan(&line)
		lines.team1LineLoss = append(lines.team1LineLoss, line)
	}

	for _, lineSlice := range lines.team1LineWin {
		q := `INSERT INTO event_winners (line_players, team, event_id, game_id)
				VALUES ($1, $2, $3, $4)`
		_, err := Db.Exec(q, lineSlice.LineTmp, lines.team1Name, lineSlice.EventID, gameID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	for _, lineSlice := range lines.team2LineLoss {
		q := `INSERT INTO event_losers (line_players, team, event_id, game_id)
				VALUES ($1, $2, $3, $4)`
		_, err := Db.Exec(q, lineSlice.LineTmp, lines.team2Name, lineSlice.EventID, gameID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	for _, lineSlice := range lines.team1LineLoss {
		q := `INSERT INTO event_losers (line_players, team, event_id, game_id)
				VALUES ($1, $2, $3, $4)`
		_, err := Db.Exec(q, lineSlice.LineTmp, lines.team1Name, lineSlice.EventID, gameID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	for _, lineSlice := range lines.team2LineWin {
		q := `INSERT INTO event_winners (line_players, team, event_id, game_id)
				VALUES ($1, $2, $3, $4)`
		_, err := Db.Exec(q, lineSlice.LineTmp, lines.team2Name, lineSlice.EventID, gameID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil

}

func compareLine(gameID string, category string) {

}
