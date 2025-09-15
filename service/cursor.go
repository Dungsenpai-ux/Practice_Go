package service

import (
	"context"
	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/Dungsenpai-ux/Practice_Go/config"
)

func GetMoviesByOffset(ctx context.Context, offset, size int) ([]model.Movie, error) {
	rows, err := config.Pool.Query(ctx, "SELECT id, title, year, genres FROM movies ORDER BY id LIMIT $1 OFFSET $2", size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, rows.Err()
}