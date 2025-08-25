package models

type StatsRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Stat struct {
	Ts string `json:"ts"`
	V  int    `json:"v"`
}

type StatsResponse struct {
	Stats []Stat `json:"stats"`
}
