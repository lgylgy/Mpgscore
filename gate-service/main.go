package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	registerHandlers()
}

func registerHandlers() {
	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/mpg", http.StatusFound))
	routes.Methods("GET").Path("/mpg").
		Handler(handler(myPlayers))

	log.Fatal(http.ListenAndServe(":8082", routes))
}

func myPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	query := r.URL.Query()
	league := query.Get("league")
	if league == "" {
		return nil, mpgErrorf("/mpg", errors.New("league parameter is missing"))
	}
	key := query.Get("key")
	if key == "" {
		return nil, mpgErrorf("/mpg", errors.New("key parameter is missing"))
	}
	result, err := getPlayers(league, key)
	if err != nil {
		return nil, mpgErrorf("/mpg", err)
	}
	return result, nil
}

type handler func(http.ResponseWriter, *http.Request) (interface{}, *mpgError)

type mpgError struct {
	Error error
	Route string
	Code  int
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, e := h(w, r)
	if e != nil {
		log.Printf("%s %d: %s", e.Route, e.Code, e.Error.Error())
		http.Error(w, e.Error.Error(), e.Code)
		return
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

func mpgErrorf(route string, err error) *mpgError {
	return &mpgError{
		Error: err,
		Route: route,
		Code:  http.StatusInternalServerError,
	}
}
