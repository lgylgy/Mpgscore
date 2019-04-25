package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	registerHandlers()
}

func registerHandlers() {
	routes := mux.NewRouter()
	routes.Handle("/", http.RedirectHandler("/mympg", http.StatusFound))
	routes.Methods("GET").Path("/mympg").
		Handler(handler(myPlayers))

	log.Fatal(http.ListenAndServe(":8080", routes))
}

func myPlayers(w http.ResponseWriter, r *http.Request) (interface{}, *mpgError) {
	query := r.URL.Query()
	league := query.Get("league")
	if league == "" {
		return nil, mpgErrorf(errors.New("missing parameter"), "invalid request")
	}
	key := query.Get("key")
	if key == "" {
		return nil, mpgErrorf(errors.New("missing parameter"), "invalid request")
	}
	result, err := getPlayers(league, key)
	if err != nil {
		return nil, mpgErrorf(err, "invalid request")
	}
	return result, nil
}

type handler func(http.ResponseWriter, *http.Request) (interface{}, *mpgError)

type mpgError struct {
	Error   error
	Message string
	Code    int
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, e := h(w, r)
	if e != nil {
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)
		http.Error(w, e.Message, e.Code)
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

func mpgErrorf(err error, format string, v ...interface{}) *mpgError {
	return &mpgError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    http.StatusInternalServerError,
	}
}
