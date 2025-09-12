package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/Dungsenpai-ux/Practice_Go/models"
	"strings"
)

func InsertMovie(ctx context.Context, movie models.Movie) (int, error) {
	var id int
	err := Pool.QueryRow(ctx, "INSERT INTO movies (title, year, genres) VALUES ($1, $2, $3) RETURNING id", movie.Title, movie.Year, movie.Genres).Scan(&id)
	return id, err
}

func GetMovieByID(ctx context.Context, id int) (models.Movie, error) {
	var movie models.Movie
	err := Pool.QueryRow(ctx, "SELECT id, title, year, genres FROM movies WHERE id = $1", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres)
	if err == pgx.ErrNoRows {
		return movie, fmt.Errorf("không tìm thấy phim")
	}
	return movie, err
}

func SearchMovies(ctx context.Context, q string, year int) ([]models.Movie, error) {
	var rows pgx.Rows
	var err error
	if year > 0 {
		rows, err = Pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1 AND year = $2", "%"+q+"%", year)
	} else {
		rows, err = Pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1", "%"+q+"%")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, rows.Err()
}

func MeasureQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
	rows, err := Pool.Query(ctx, "EXPLAIN ANALYZE "+query, args...)
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