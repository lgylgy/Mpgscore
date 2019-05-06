package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mpgscore/api"
	"net/http"
	"os"
	"strconv"
)

var mpgPort int = retrievePort("PLAYERDB")
var dbPort int = retrievePort("MPGDB")

func main() {
	registerHandlers()
}

func retrievePort(value string) int {
	key := os.Getenv(value)
	if len(key) == 0 {
		return 0
	}
	port, err := strconv.Atoi(key)
	if err != nil {
		return 0
	}
	return port
}

func registerHandlers() {
	port := retrievePort("GATEDB")
	if port == 0 {
		log.Fatal("$GATEDB variable is not present")
	}

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
