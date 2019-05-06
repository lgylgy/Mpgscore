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

func main() {
	registerHandlers()
}

func registerHandlers() {
	key := os.Getenv("GATEDB")
	if len(key) == 0 {
		log.Fatal("$GATEDB variable is not present")
	}
	port, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal(err)
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
