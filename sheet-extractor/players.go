package main

import (
	"fmt"
	"sync"

	"github.com/lgylgy/mpgscore/api"
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

func extractPlayer(row *sheets.RowData, team string, matchs int) *api.DbPlayer {
	player := &api.DbPlayer{}
	player.Team = api.NormalizeString(team)
	for c, data := range row.Values {
		if c == nameRow {
			player.Name = api.NormalizeString(data.FormattedValue)
		} else if c >= matchsCols && c < matchsCols+matchs {
			player.Grades = append(player.Grades, data.FormattedValue)
		}
	}
	if player.Name == "" {
		return nil
	}
	return player
}

func validMetaData(team string, matchCount int) error {
	if team == "" {
		return fmt.Errorf("invalid name team in %v", team)
	}
	if matchCount == 0 {
		return fmt.Errorf("no match is filled in %v", team)
	}
	return nil
}

func extractTeamPlayers(sheet *sheets.Sheet, teamID string) ([]*api.DbPlayer, error) {
	if len(sheet.Data) == 0 {
		return nil, fmt.Errorf("any data in the grid %v", teamID)
	}
	var matchCount int
	players := []*api.DbPlayer{}
	for r, row := range sheet.Data[0].RowData {
		if r == matchsRow {
			matchCount = extractMatchCount(row)
		} else if r >= startingPlayer {
			err := validMetaData(teamID, matchCount)
			if err != nil {
				return nil, err
			}
			player := extractPlayer(row, teamID, matchCount)
			if player == nil {
				break
			}
			players = append(players, player)
		}
	}
	return players, nil
}

func extractPlayers(tabs []*sheets.Sheet, jobs uint, team string) ([]*api.DbPlayer, error) {

	if len(tabs) < 2 {
		return nil, fmt.Errorf("Invalid sheet count")
	}

	players := []*api.DbPlayer{}
	pending := make(chan *sheets.Sheet)
	results := make(chan []*api.DbPlayer)
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
				title := v.Properties.Title
				if team == "" || title == team {
					fmt.Printf("Extract %v sheet..\n", title)
					players, err := extractTeamPlayers(v, title)
					if err != nil {
						break
					}
					results <- players
				}
			}
		}()
	}

	for result := range results {
		players = append(players, result...)
	}

	return players, nil
}
