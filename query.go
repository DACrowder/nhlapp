package main

import (
	"fmt"
	"strconv"
	"strings"
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
	DataCount int    `db:"data_count" json:"linupCount"`
	LineTmp   string `db:"line_array"`
	Line      []int  `json:"lineup"`
}

type Lines struct {
	Team1Name string     `json:"Team1Name"`
	Team1Line []LineData `json:"Team1LineData"`
	Team2Name string     `json:"Team2Name"`
	Team2Line []LineData `json:"Team2LineData"`
}

//not safe use are your own risk
func parseLine(lines Lines) {
	for i, line := range lines.Team1Line {
		lineTmp := strings.Split(line.LineTmp[1:len(line.LineTmp)-2], ",")
		for _, p := range lineTmp {
			pInt, err := strconv.Atoi(p)
			if err != nil {
				return
			}
			lines.Team1Line[i].Line = append(lines.Team1Line[i].Line, pInt)
		}
	}
	for i, line := range lines.Team2Line {
		lineTmp := strings.Split(line.LineTmp[1:len(line.LineTmp)-2], ",")
		for _, p := range lineTmp {
			pInt, err := strconv.Atoi(p)
			if err != nil {
				return
			}
			lines.Team2Line[i].Line = append(lines.Team2Line[i].Line, pInt)
		}
	}
}

func getLineShots(gameID string) (*Lines, error) {
	q := `select distinct player1_team from event where game_id = $1 and player1_team != '';`
	//get both teams
	rows, err := Db.Queryx(q, gameID)
	if err != nil {
		return nil, err
	}

	lines := Lines{}

	if rows.Next() {
		if err = rows.Scan(&lines.Team1Name); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	if rows.Next() {
		if err = rows.Scan(&lines.Team2Name); err != nil {
			return nil, err
		}
	} else {
		return nil, err
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
	rows, err = Db.Queryx(q, gameID, lines.Team1Name)
	if err != nil {
		return nil, err
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

	fmt.Println("------ TEAM 1 -------")
	for _, lineSlice := range lines.Team1Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}
	fmt.Println("------ TEAM 2 -------")

	for _, lineSlice := range lines.Team2Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}

	parseLine(lines)

	return &lines, nil
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
						where game_id = $1 and event_type = $2 and player1_team = $3
                    ) as sub_q
                where roster.game_id = $1 and roster.event_id = sub_q.event_id and roster.team = sub_q.team
                order by roster.event_id asc) as q
        group by q.event_id) as q2
    group by q2.players order by count(*)`

	/* get first team */
	rows, err = Db.Queryx(q, gameID, category, lines.Team1Name)
	if err != nil {
		return err
	}

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.Team1Line = append(lines.Team1Line, line)
	}

	/* get second game */
	rows, err = Db.Queryx(q, gameID, category, lines.Team2Name)

	for rows.Next() {
		line := LineData{}
		rows.StructScan(&line)
		lines.Team2Line = append(lines.Team2Line, line)
	}
	fmt.Printf("\t------ %s [%s] -------\n", lines.Team1Name, category)
	for _, lineSlice := range lines.Team1Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}
	fmt.Printf("\t------ %s [%s] -------\n", lines.Team2Name, category)
	for _, lineSlice := range lines.Team2Line {
		fmt.Printf("%d  %s\n", lineSlice.DataCount, lineSlice.LineTmp)
	}

	return nil
}
