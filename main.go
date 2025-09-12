package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"log"
	"github.com/Dungsenpai-ux/Practice_Go/db"
	"github.com/Dungsenpai-ux/Practice_Go/handlers"
	"github.com/Dungsenpai-ux/Practice_Go/models"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations() error {
	m, err := migrate.New("file://migrations", "postgres://postgres:123@localhost:5432/movies_db?sslmode=disable")
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func seedData() error {
	file, err := os.Open("data/movies.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Read() // Bỏ qua header

	ctx := context.Background()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`(.*) \((\d{4})\)`)
		matches := re.FindStringSubmatch(record[1])
		title := record[1]
		year := 0
		if len(matches) == 3 {
			title = matches[1]
			year, _ = strconv.Atoi(matches[2])
		}

		movie := models.Movie{
			Title:  title,
			Year:   year,
			Genres: strings.Replace(record[2], "|", ", ", -1),
		}
		_, err = db.InsertMovie(ctx, movie)
		if err != nil {
			log.Printf("Bỏ qua: %v", err)
		}
	}
	return nil
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := runMigrations(); err != nil {
		log.Fatal(err)
	}

	if err := seedData(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("POST /movies", handlers.CreateMovie)
	http.HandleFunc("GET /movies/", handlers.GetMovie)
	http.HandleFunc("GET /movies/search", handlers.SearchMovies)

	log.Println("Máy chủ khởi động tại :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}