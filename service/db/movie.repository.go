package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertMovie(ctx context.Context, pool *pgxpool.Pool, movie model.Movie) (int, error) {
	var id int
	err := pool.QueryRow(ctx, "INSERT INTO movies (title, year, genres) VALUES ($1,$2,$3) RETURNING id", movie.Title, movie.Year, movie.Genres).Scan(&id)
	return id, err
}

func GetMovieByID(ctx context.Context, pool *pgxpool.Pool, id int) (model.Movie, error) {
	var movie model.Movie
	err := pool.QueryRow(ctx, "SELECT id, title, year, genres FROM movies WHERE id=$1", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres)
	if err == pgx.ErrNoRows {
		return movie, fmt.Errorf("không tìm thấy phim")
	}
	return movie, err
}

func SearchMovies(ctx context.Context, pool *pgxpool.Pool, q string, year int) ([]model.Movie, error) {
	var rows pgx.Rows
	var err error
	if year > 0 {
		rows, err = pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1 AND year=$2", "%"+q+"%", year)
	} else {
		rows, err = pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1", "%"+q+"%")
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

func MeasureQuery(ctx context.Context, pool *pgxpool.Pool, query string, args ...interface{}) (string, error) {
	rows, err := pool.Query(ctx, "EXPLAIN ANALYZE "+query, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var sb strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return "", err
		}
		sb.WriteString(line + "\n")
	}
	return sb.String(), nil
}
