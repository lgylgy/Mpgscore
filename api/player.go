package api

type Player struct {
	Name   string   `json:"name"`
	Team   string   `json:"team"`
	Grades []string `json:"grades"`
}
