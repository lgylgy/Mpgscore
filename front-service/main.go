package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"

	api "github.com/lgylgy/mpgscore/api"
)

var (
	gateAddr  = ""
	templates = template.Must(template.New("").ParseGlob("templates/*.html"))
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

	gateAddr = api.GetEnv("GATE_SERVICE_ADDR")

	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/mpg", http.StatusFound))
	routes.HandleFunc("/mpg", myPlayers)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routes))
}

func myPlayers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	result, err := getPlayers(query.Get("league"), query.Get("key"))
	if err != nil {
		log.Println(err)
	}

	err = templates.ExecuteTemplate(w, "players", map[string]interface{}{
		"players": result,
	})
	if err != nil {
		log.Println(err)
	}
}
