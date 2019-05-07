package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mpgscore/api"
	"net/http"
	"net/url"
)

func getMyTeam(league, key string) ([]*api.MpgPlayer, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://localhost:%d/mympg?league=%s&key=%s", mpgPort, league, key)
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

	var info api.MpgData
	err = json.Unmarshal(body, &info.Players)
	if err != nil {
		return nil, err
	}
	return info.Players, nil
}

func getMyPlayer(firstname, lastname string) (*api.DbPlayer, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://localhost:%d/player?firstname=%s&lastname=%s",
		dbPort, url.QueryEscape(firstname), url.QueryEscape(lastname))
	req, err := http.NewRequest("GET", url, nil)
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

	var info api.DbPlayer
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func getPlayers(league, key string) ([]*api.DbPlayer, error) {

	myTeam, err := getMyTeam(league, key)
	if err != nil {
		return nil, err
	}

	result := []*api.DbPlayer{}
	for _, v := range myTeam {
		player, err := getMyPlayer(api.NormalizeString(v.Firstname),
			api.NormalizeString(v.Lastname))
		if err != nil {
			continue
		}
		result = append(result, player)
	}

	return result, nil
}
