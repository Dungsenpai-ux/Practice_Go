package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func InsertMovie(ctx context.Context, pool *pgxpool.Pool, movie model.Movie) (int, error) {
	ctx, span := otel.Tracer("repo.movies").Start(ctx, "InsertMovie")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.system", "postgres"),
		attribute.String("db.operation", "INSERT"),
		attribute.String("movie.title", movie.Title),
		attribute.Int("movie.year", movie.Year),
	)
	var id int
	start := time.Now()
	err := pool.QueryRow(ctx, "INSERT INTO movies (title, year, genres) VALUES ($1,$2,$3) RETURNING id", movie.Title, movie.Year, movie.Genres).Scan(&id)
	span.SetAttributes(attribute.Int64("db.duration_ms", time.Since(start).Milliseconds()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return id, err
}

func GetMovieByID(ctx context.Context, pool *pgxpool.Pool, id int) (model.Movie, error) {
	ctx, span := otel.Tracer("repo.movies").Start(ctx, "GetMovieByID")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.system", "postgres"),
		attribute.String("db.operation", "SELECT"),
		attribute.Int("movie.id", id),
	)
	var movie model.Movie
	start := time.Now()
	err := pool.QueryRow(ctx, "SELECT id, title, year, genres FROM movies WHERE id=$1", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres)
	span.SetAttributes(attribute.Int64("db.duration_ms", time.Since(start).Milliseconds()))
	if err == pgx.ErrNoRows {
		err = fmt.Errorf("không tìm thấy phim")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return movie, err
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return movie, err
}

func SearchMovies(ctx context.Context, pool *pgxpool.Pool, q string, year int) ([]model.Movie, error) {
	ctx, span := otel.Tracer("repo.movies").Start(ctx, "SearchMovies")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.system", "postgres"),
		attribute.String("db.operation", "SELECT"),
		attribute.String("search.query", q),
		attribute.Int("search.year", year),
	)
	var rows pgx.Rows
	var err error
	queryStart := time.Now()
	if year > 0 {
		rows, err = pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1 AND year=$2", "%"+q+"%", year)
	} else {
		rows, err = pool.Query(ctx, "SELECT id, title, year, genres FROM movies WHERE title ILIKE $1", "%"+q+"%")
	}
	span.SetAttributes(attribute.Int64("db.duration_ms", time.Since(queryStart).Milliseconds()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer rows.Close()
	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genres); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return movies, nil
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
