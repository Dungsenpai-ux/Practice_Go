package service

import (
	"context"

	"github.com/Dungsenpai-ux/Practice_Go/model"
	repo "github.com/Dungsenpai-ux/Practice_Go/service/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MovieService holds business logic for movies independent from transport.
type MovieService struct{ pool *pgxpool.Pool }

func NewMovieService(pool *pgxpool.Pool) *MovieService { return &MovieService{pool: pool} }

func (s *MovieService) InsertMovie(ctx context.Context, movie model.Movie) (int, error) {
	return repo.InsertMovie(ctx, s.pool, movie)
}

func (s *MovieService) GetMovieByID(ctx context.Context, id int) (model.Movie, error) {
	return repo.GetMovieByID(ctx, s.pool, id)
}

func (s *MovieService) SearchMovies(ctx context.Context, q string, year int) ([]model.Movie, error) {
	return repo.SearchMovies(ctx, s.pool, q, year)
}

func (s *MovieService) MeasureQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
	return repo.MeasureQuery(ctx, s.pool, query, args...)
}
