package model

type Movie struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Year   int    `json:"year,omitempty"`
	Genres string `json:"genres,omitempty"`
}
