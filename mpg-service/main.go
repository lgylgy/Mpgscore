package main

import (
	"errors"
	"github.com/gorilla/mux"
	"log"
	"mpgscore/api"
	"net/http"
)

func main() {
	registerHandlers()
}

func registerHandlers() {
	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/mympg", http.StatusFound))
	routes.Methods("GET").Path("/mympg").
		Handler(api.Handler(myPlayers))

	log.Fatal(http.ListenAndServe(":8080", routes))
}

func myPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	query := r.URL.Query()
	league := query.Get("league")
	if league == "" {
		return nil, api.MpgErrorf("/mympg", errors.New("missing parameter"))
	}
	key := query.Get("key")
	if key == "" {
		return nil, api.MpgErrorf("/mympg", errors.New("missing parameter"))
	}
	result, err := getPlayers(league, key)
	if err != nil {
		return nil, api.MpgErrorf("/mympg", err)
	}
	return result, nil
}
