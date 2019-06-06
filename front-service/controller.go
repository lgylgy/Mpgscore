package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	api "github.com/lgylgy/mpgscore/api"
)

func getPlayers(league, key string) ([]*api.DbPlayer, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s/mpg?league=%s&key=%s", gateAddr, league, key)
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

	result := []*api.DbPlayer{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
