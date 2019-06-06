package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	api "github.com/lgylgy/mpgscore/api"
)

var (
	mpgAddr = ""
	dbAddr  = ""
)

func main() {
	registerHandlers()
}

func registerHandlers() {
	key := api.GetEnv("PORT")
	port, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal(err)
	}

	mpgAddr = api.GetEnv("MPG_SERVICE_ADDR")
	dbAddr = api.GetEnv("DB_SERVICE_ADDR")

	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/mpg", http.StatusFound))
	routes.Methods("GET").Path("/mpg").
		Handler(api.Handler(myPlayers))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routes))
}

func myPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	query := r.URL.Query()
	league := query.Get("league")
	if league == "" {
		return nil, api.MpgErrorf("/mpg", errors.New("league parameter is missing"))
	}
	key := query.Get("key")
	if key == "" {
		return nil, api.MpgErrorf("/mpg", errors.New("key parameter is missing"))
	}
	result, err := getPlayers(league, key)
	if err != nil {
		return nil, api.MpgErrorf("/mpg", err)
	}
	return result, nil
}
