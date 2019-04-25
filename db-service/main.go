package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Player struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Team   string   `json:"team"`
	Grades []string `json:"grades"`
}

var controller *Controller

func main() {

	controller = NewController()
	err := controller.Connect("$MONGODB", "$MONGODB", "$MONGODB")
	if err != nil {
		log.Fatal(err)
	}
	defer controller.Close()

	registerHandlers()
}

func registerHandlers() {
	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/players", http.StatusFound))
	routes.Methods("GET").Path("/players").
		Handler(handler(listPlayers))
	routes.Methods("GET").Path("/players/{id}").
		Handler(handler(getPlayerById))
	routes.Methods("GET").Path("/player").
		Handler(handler(getPlayer))
	routes.Methods("GET").Path("/teams/{id}").
		Handler(handler(listTeamPlayers))
	routes.Methods("PUT").Path("/players/{id}").
		Handler(handler(updatePlayer))
	routes.Methods("POST").Path("/players").
		Handler(handler(createPlayer))

	log.Fatal(http.ListenAndServe(":8080", routes))
}

func createPlayer(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	defer r.Body.Close()

	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		return nil, mpgErrorf("POST /players", err)
	}
	_, err := controller.AddPlayer(&player)
	if err != nil {
		return nil, mpgErrorf("POST /players", err)
	}
	return player, nil
}

func listPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	players, err := controller.ListPlayers()
	if err != nil {
		return nil, mpgErrorf("GET /players", err)
	}
	return players, nil
}

func listTeamPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	params := mux.Vars(r)
	players, err := controller.ListTeamPlayers(params["id"])
	if err != nil {
		return nil, mpgErrorf("GET /teams", err)
	}
	return players, nil
}

func getPlayerById(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	params := mux.Vars(r)
	player, err := controller.GetPlayerById(params["id"])
	if err != nil {
		return nil, mpgErrorf("GET /player", err)
	}
	return player, nil
}

func getPlayer(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	query := r.URL.Query()
	firstname := query.Get("firstname")
	if firstname == "" {
		return nil, mpgErrorf("GET /player", errors.New("missing parameter"))
	}
	lastname := query.Get("lastname")
	if lastname == "" {
		return nil, mpgErrorf("GET /player", errors.New("missing parameter"))
	}

	player, err := controller.GetPlayer(firstname, lastname)
	if err != nil {
		return nil, mpgErrorf("GET /player", err)
	}
	return player, nil
}

func updatePlayer(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		return nil, mpgErrorf("PUT /players", err)
	}
	params := mux.Vars(r)
	player.ID = params["id"]
	updated, err := controller.UpdatePlayer(&player)
	if err != nil {
		return nil, mpgErrorf("PUT /players", err)
	}
	return updated, nil
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
