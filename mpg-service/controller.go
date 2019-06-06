package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	api "github.com/lgylgy/mpgscore/api"
)

type PlayersInfo struct {
	Data *api.MpgData `json:"data"`
}

type TeamData struct {
	Id      string           `json:"id"`
	Players []*api.MpgPlayer `json:"players"`
}

type TeamsInfo struct {
	CurrentTeam string               `json:"current_mpg_team"`
	Teams       map[string]*TeamData `json:"teams"`
}

func getPlayers(league, key, remote, route string) ([]*api.MpgPlayer, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://%s/%s/%s", remote, league, route)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", key)
	rpy, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rpy.Body.Close()

	body, err := ioutil.ReadAll(rpy.Body)
	if err != nil {
		return nil, err
	}

	if route == "coach" {
		return getCoachPlayers(body)
	}
	return getTeamsPlayers(body)
}

func getCoachPlayers(body []byte) ([]*api.MpgPlayer, error) {
	var info PlayersInfo
	err := json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	if info.Data != nil {
		return info.Data.Players, nil
	}
	return []*api.MpgPlayer{}, nil
}

func getTeamsPlayers(body []byte) ([]*api.MpgPlayer, error) {
	var info TeamsInfo
	err := json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	for k, v := range info.Teams {
		if k == info.CurrentTeam {
			return v.Players, nil
		}
	}
	return []*api.MpgPlayer{}, nil
}
