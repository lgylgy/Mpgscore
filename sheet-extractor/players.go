package main

import (
	"Mpgscore/api"
	"fmt"
	"sync"

	"google.golang.org/api/sheets/v4"
)

const (
	startingSheet  = 1
	startingPlayer = 7
	teamCol        = 0
	teamRow        = 0
	matchsRow      = 6
	matchsCols     = 7
	nameRow        = 2
)

func extractMatchCount(row *sheets.RowData) int {
	count := 0
	for c, data := range row.Values {
		if c >= matchsCols {
			if data.FormattedValue != "" {
				count++
			} else {
				return count
			}
		}
	}
	return count
}

func extractTeamName(row *sheets.RowData) string {
	for c, data := range row.Values {
		if c == teamCol {
			return data.FormattedValue
		}
	}
	return ""
}

func extractPlayer(row *sheets.RowData, team string, matchs int) *api.Player {
	player := &api.Player{}
	player.Team = team
	for c, data := range row.Values {
		if c == nameRow {
			player.Name = data.FormattedValue
		} else if c >= matchsCols && c < matchsCols+matchs {
			player.Grades = append(player.Grades, data.FormattedValue)
		}
	}
	if player.Name == "" {
		return nil
	}
	return player
}

func validMetaData(team, sheetTitle string, matchCount int) error {
	if team == "" {
		return fmt.Errorf("invalid name team in %v", sheetTitle)
	}
	if matchCount == 0 {
		return fmt.Errorf("no match is filled in %v", sheetTitle)
	}
	return nil
}

func extractTeamPlayers(sheet *sheets.Sheet) ([]*api.Player, error) {
	sheetTitle := sheet.Properties.Title
	fmt.Printf("Extract %v sheet..\n", sheetTitle)
	if len(sheet.Data) == 0 {
		return nil, fmt.Errorf("any data in the grid %v", sheetTitle)
	}
	var team string
	var matchCount int
	players := []*api.Player{}
	for r, row := range sheet.Data[0].RowData {
		if r == teamRow {
			team = extractTeamName(row)
		} else if r == matchsRow {
			matchCount = extractMatchCount(row)
		} else if r >= startingPlayer {
			err := validMetaData(team, sheetTitle, matchCount)
			if err != nil {
				return nil, err
			}
			player := extractPlayer(row, team, matchCount)
			if player == nil {
				break
			}
			players = append(players, player)
		}
	}
	return players, nil
}

func extractPlayers(tabs []*sheets.Sheet, jobs uint) ([]*api.Player, error) {

	if len(tabs) < 2 {
		return nil, fmt.Errorf("Invalid sheet count")
	}

	players := []*api.Player{}
	pending := make(chan *sheets.Sheet)
	results := make(chan []*api.Player)
	stop := make(chan bool)
	running := &sync.WaitGroup{}
	running.Add(int(jobs))

	tabs = tabs[1:]
	go func() {
	Outer:
		for _, s := range tabs {
			select {
			case pending <- s:
			case <-stop:
				break Outer
			}
		}
		close(pending)
		running.Wait()
		close(results)
	}()

	for i := uint(0); i < jobs; i++ {
		go func() {
			defer running.Done()
			for v := range pending {
				players, err := extractTeamPlayers(v)
				if err != nil {
					break
				}
				results <- players
			}
		}()
	}

	for result := range results {
		players = append(players, result...)
	}

	return players, nil
}
