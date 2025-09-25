package db

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SeedData loads data/movies.csv and inserts rows into movies table if table is empty.
func SeedData(pool *pgxpool.Pool) error {
	// Check if table already has data
	var count int
	if err := pool.QueryRow(ctxBackground(), "SELECT COUNT(1) FROM movies").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	f, err := os.Open("data/movies.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	// read header
	if _, err := reader.Read(); err != nil {
		return err
	}

	// Insert in batches
	const batchSize = 1000
	type row struct {
		title  string
		year   int
		genres string
	}
	batch := make([]row, 0, batchSize)

	flush := func(rows []row) error {
		if len(rows) == 0 {
			return nil
		}
		// Build multi-values insert
		var sb strings.Builder
		args := make([]interface{}, 0, len(rows)*3)
		sb.WriteString("INSERT INTO movies (title, year, genres) VALUES ")
		for i, r := range rows {
			if i > 0 {
				sb.WriteString(",")
			}
			idx := i * 3
			sb.WriteString(fmt.Sprintf("($%d,$%d,$%d)", idx+1, idx+2, idx+3))
			args = append(args, r.title, r.year, r.genres)
		}
		_, err := pool.Exec(ctxBackground(), sb.String(), args...)
		return err
	}

	for {
		rec, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if len(rec) < 3 {
			continue
		}
		// rec[0]=movieId, rec[1]=title (may include year in parentheses), rec[2]=genres
		title := rec[1]
		year := parseYearFromTitle(title)
		genres := rec[2]
		batch = append(batch, row{title: title, year: year, genres: genres})
		if len(batch) >= batchSize {
			if err := flush(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}
	if err := flush(batch); err != nil {
		return err
	}
	return nil
}

func parseYearFromTitle(title string) int {
	// Expect titles like "Toy Story (1995)"; extract digits inside last parentheses
	open := strings.LastIndex(title, "(")
	close := strings.LastIndex(title, ")")
	if open >= 0 && close > open+1 {
		if y, err := strconv.Atoi(strings.TrimSpace(title[open+1 : close])); err == nil {
			return y
		}
	}
	return 0
}

// ctxBackground isolates dependency from callers
func ctxBackground() context.Context { return context.Background() }
