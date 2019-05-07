package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mpgscore/api"
	"net/http"
	"os"
	"strconv"
)

var controller *Controller

func main() {
	key := os.Getenv("MPGDB")
	if len(key) == 0 {
		log.Fatal("$MPGDB variable is not present")
	}
	port, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal(err)
	}

	mongoDB := os.Getenv("MONGODB")
	if len(mongoDB) == 0 {
		log.Fatal("$MONGODB variable is not present")
	}

	controller = NewController()
	err = controller.Connect(mongoDB, "mpg", "mpg", false)
	if err != nil {
		log.Fatal(err)
	}
	defer controller.Close()
	log.Println("Connection MongoDB succeed!")

	registerHandlers(port)
}

func registerHandlers(port int) {
	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/players", http.StatusFound))
	routes.Methods("GET").Path("/players").
		Handler(api.Handler(listPlayers))
	routes.Methods("GET").Path("/players/{id}").
		Handler(api.Handler(getPlayerById))
	routes.Methods("GET").Path("/player").
		Handler(api.Handler(getPlayer))
	routes.Methods("GET").Path("/teams/{id}").
		Handler(api.Handler(listTeamPlayers))
	routes.Methods("PUT").Path("/players/{id}").
		Handler(api.Handler(updatePlayer))
	routes.Methods("POST").Path("/players").
		Handler(api.Handler(createPlayer))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routes))
}

func createPlayer(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	defer r.Body.Close()

	var player api.DbPlayer
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		return nil, api.MpgErrorf("POST /players", err)
	}
	_, err := controller.AddPlayer(&player)
	if err != nil {
		return nil, api.MpgErrorf("POST /players", err)
	}
	return player, nil
}

func listPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	players, err := controller.ListPlayers()
	if err != nil {
		return nil, api.MpgErrorf("GET /players", err)
	}
	return players, nil
}

func listTeamPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	params := mux.Vars(r)
	players, err := controller.ListTeamPlayers(params["id"])
	if err != nil {
		return nil, api.MpgErrorf("GET /teams", err)
	}
	return players, nil
}

func getPlayerById(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	params := mux.Vars(r)
	player, err := controller.GetPlayerById(params["id"])
	if err != nil {
		return nil, api.MpgErrorf("GET /player", err)
	}
	return player, nil
}

func getPlayer(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	query := r.URL.Query()
	firstname := query.Get("firstname")
	lastname := query.Get("lastname")
	if firstname == "" && lastname == "" {
		return nil, api.MpgErrorf("GET /player", errors.New("missing parameter"))
	}
	player, err := controller.GetPlayer(firstname, lastname)
	if err != nil {
		return nil, api.MpgErrorf("GET /player", err)
	}
	return player, nil
}

func updatePlayer(w http.ResponseWriter, r *http.Request) (interface{}, *api.MpgError) {
	var player api.DbPlayer
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		return nil, api.MpgErrorf("PUT /players", err)
	}
	params := mux.Vars(r)
	player.ID = params["id"]
	updated, err := controller.UpdatePlayer(&player)
	if err != nil {
		return nil, api.MpgErrorf("PUT /players", err)
	}
	return updated, nil
}
