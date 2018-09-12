package api

type Player struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Team   string   `json:"team"`
	Grades []string `json:"grades"`
}
