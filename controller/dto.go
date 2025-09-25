package controller

import "github.com/Dungsenpai-ux/Practice_Go/model"

// MovieCreateRequest DTO for creating a movie (controller layer only)
type MovieCreateRequest struct {
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Genres string `json:"genres"`
}

// MovieResponse DTO returned to clients
type MovieResponse struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Year   int    `json:"year,omitempty"`
	Genres string `json:"genres,omitempty"`
}

func (r MovieCreateRequest) ToModel() model.Movie {
	return model.Movie{Title: r.Title, Year: r.Year, Genres: r.Genres}
}

func FromModel(m model.Movie) MovieResponse {
	return MovieResponse{ID: m.ID, Title: m.Title, Year: m.Year, Genres: m.Genres}
}

func FromModelSlice(ms []model.Movie) []MovieResponse {
	out := make([]MovieResponse, 0, len(ms))
	for _, m := range ms {
		out = append(out, FromModel(m))
	}
	return out
}

// HealthResponseDTO used by health endpoint
type HealthResponseDTO struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}
