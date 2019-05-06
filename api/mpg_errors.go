package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) (interface{}, *MpgError)

type MpgError struct {
	Error error
	Route string
	Code  int
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func MpgErrorf(route string, err error) *MpgError {
	return &MpgError{
		Error: err,
		Route: route,
		Code:  http.StatusInternalServerError,
	}
}
