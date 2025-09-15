package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/jackc/pgx/v5"
)

func InsertMovie(ctx context.Context, movie model.Movie) (int, error) {
	var id int
	err := config.Pool.QueryRow(ctx, "INSERT INTO movies (title, year, genres) VALUES ($1, $2, $3) RETURNING id", movie.Title, movie.Year, movie.Genres).Scan(&id)
	return id, err
}

func GetMovieByID(ctx context.Context, id int) (model.Movie, error) {
	var movie model.Movie
	err := config.Pool.QueryRow(ctx, "SELECT id, title, year, genres FROM movies WHERE id = $1", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres)
	if err == pgx.ErrNoRows {
		return movie, fmt.Errorf("không tìm thấy phim")
	}
	return movie, err
}

func SearchMovies(ctx context.Context, q string, year int) ([]model.Movie, error) {
	var rows pgx.Rows
	var err error
	if year > 0 {
		rows, err = config.Pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1 AND year = $2", "%"+q+"%", year)
	} else {
		rows, err = config.Pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1", "%"+q+"%")
	}
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

func GetMoviesByCursor(ctx context.Context, cursor string, size int) ([]model.Movie, string, error) {
	var rows pgx.Rows
	var err error
	var nextCursor string

	if cursor == "" {
		rows, err = config.Pool.Query(ctx, "SELECT id, title, year, genres FROM movies ORDER BY id LIMIT $1", size)
	} else {
		rows, err = config.Pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE id > $1 ORDER BY id LIMIT $2", cursor, size)
	}

	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres); err != nil {
			return nil, "", err
		}
		movies = append(movies, movie)
		nextCursor = strconv.Itoa(movie.ID)
	}
	return movies, nextCursor, rows.Err()
}

func LogQueryPerformance(ctx context.Context, query string, args ...interface{}) {
	start := time.Now()
	_, err := config.Pool.Query(ctx, query, args...)
	duration := time.Since(start)
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
	} else {
		fmt.Printf("Query executed in: %v\n", duration)
	}
}
