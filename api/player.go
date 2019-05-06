package api

type DbPlayer struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Team   string   `json:"team"`
	Grades []string `json:"grades"`
}

type MpgPlayer struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type MpgData struct {
	Players []*MpgPlayer `json:"players"`
}
