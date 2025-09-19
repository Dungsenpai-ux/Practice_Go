package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/jackc/pgx/v5"
)

func InsertMovie(ctx context.Context, movie model.Movie) (int, error) {
	var id int
	timer := config.ObserveDBQueryDuration()
	err := config.Pool.QueryRow(ctx, "INSERT INTO movies (title, year, genres) VALUES ($1, $2, $3) RETURNING id", movie.Title, movie.Year, movie.Genres).Scan(&id)
	timer()
	return id, err
}

func GetMovieByID(ctx context.Context, id int) (model.Movie, error) {
	var movie model.Movie
	timer := config.ObserveDBQueryDuration()
	err := config.Pool.QueryRow(ctx, "SELECT id, title, year, genres FROM movies WHERE id = $1", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres)
	timer()
	if err == pgx.ErrNoRows {
		return movie, fmt.Errorf("không tìm thấy phim")
	}
	return movie, err
}

func SearchMovies(ctx context.Context, q string, year int) ([]model.Movie, error) {
	var rows pgx.Rows
	var err error
	timer := config.ObserveDBQueryDuration()
	if year > 0 {
		rows, err = config.Pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1 AND year = $2", "%"+q+"%", year)
	} else {
		rows, err = config.Pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1", "%"+q+"%")
	}
	timer()
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

func MeasureQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
	rows, err := config.Pool.Query(ctx, "EXPLAIN ANALYZE "+query, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var explain strings.Builder
	for rows.Next() {
		var line string
		rows.Scan(&line)
		explain.WriteString(line + "\n")
	}
	return explain.String(), nil
}
