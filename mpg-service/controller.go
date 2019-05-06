package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mpgscore/api"
	"net/http"
)

type PlayersInfo struct {
	Data *api.MpgData `json:"data"`
}

func getPlayers(league, key string) ([]*api.MpgPlayer, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.monpetitgazon.com/league/%s/coach", league)
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

	var info PlayersInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	if info.Data != nil {
		return info.Data.Players, nil
	}
	return []*api.MpgPlayer{}, nil
}
