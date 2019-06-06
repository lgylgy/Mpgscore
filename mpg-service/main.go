package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"

	api "github.com/lgylgy/mpgscore/api"
)

func main() {
	registerHandlers()
}

func registerHandlers() {
	key := os.Getenv("PLAYERDB")
	if len(key) == 0 {
		log.Fatal("$PLAYERDB variable is not present")
	}
	port, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("port: %v\n", port)

	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/mympg", http.StatusFound))
	routes.Methods("GET").Path("/mympg").
		Handler(api.Handler(myPlayers))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routes))
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
