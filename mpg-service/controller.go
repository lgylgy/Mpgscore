package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getPlayers(league, key string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.monpetitgazon.com/league/%s/coach", league)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", key)
	rpy, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rpy.Body.Close()
	body, err := ioutil.ReadAll(rpy.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
