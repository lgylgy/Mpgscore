package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type mPlayer struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type mData struct {
	Players []*mPlayer `json:"players"`
}

type xPlayer struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Team   string   `json:"team"`
	Grades []string `json:"grades"`
}

func getMyTeam(league, key string) ([]*mPlayer, error) {
	client := &http.Client{}

	url := fmt.Sprintf("$MYMPG", league, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	rpy, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rpy.Body.Close()

	body, err := ioutil.ReadAll(rpy.Body)
	if err != nil {
		return nil, err
	}

	var info mData
	err = json.Unmarshal(body, &info.Players)
	if err != nil {
		return nil, err
	}
	return info.Players, nil
}

func getMyPlayer(firstname, lastname string) (*xPlayer, error) {
	client := &http.Client{}

	adr := fmt.Sprintf("$MYPLAYER",
		url.QueryEscape(firstname), url.QueryEscape(lastname))
	req, err := http.NewRequest("GET", adr, nil)
	if err != nil {
		return nil, err
	}

	rpy, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rpy.Body.Close()

	if rpy.StatusCode != http.StatusOK {
		return nil, errors.New("invalid request")
	}

	body, err := ioutil.ReadAll(rpy.Body)
	if err != nil {
		return nil, err
	}

	var info xPlayer
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func getPlayers(league, key string) ([]*xPlayer, error) {

	myTeam, err := getMyTeam(league, key)
	if err != nil {
		return nil, err
	}

	result := []*xPlayer{}
	for _, v := range myTeam {
		player, err := getMyPlayer(v.Firstname, v.Lastname)
		if err != nil {
			continue
		}
		result = append(result, player)
	}

	return result, nil
}
