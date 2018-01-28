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

type Lines struct {
	team1Name string
	team1Line []LineData
	team2Name string
	team2Line []LineData
}

func getLineShots(gameID string) error {
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
                        where game_id = $1 and (event_type = 'SHOT' or event_type = 'GOAL') and player1_team = $2
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

	fmt.Println("------ TEAM 1 -------")
	for _, lineSlice := range lines.team1Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}
	fmt.Println("------ TEAM 2 -------")

	for _, lineSlice := range lines.team2Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}

	return nil
}

func getWildCard(gameID string, category string) error {
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
						where game_id = $1 and event_type = $2 and player1_team = $3
                    ) as sub_q
                where roster.game_id = $1 and roster.event_id = sub_q.event_id and roster.team = sub_q.team
                order by roster.event_id asc) as q
        group by q.event_id) as q2
    group by q2.players order by count(*)`

	/* get first team */
	rows, err = Db.Queryx(q, gameID, category, lines.team1Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.team1Line = append(lines.team1Line, line)
	}

	/* get second game */
	rows, err = Db.Queryx(q, gameID, category, lines.team2Name)

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.team2Line = append(lines.team2Line, line)
	}
	fmt.Printf("\t------ %s [%s] -------\n", lines.team1Name, category)
	for _, lineSlice := range lines.team1Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}
	fmt.Printf("\t------ %s [%s] -------\n", lines.team2Name, category)
	for _, lineSlice := range lines.team2Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}

	return nil
}
