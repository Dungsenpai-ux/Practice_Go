package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"log"
	// "github.com/gin-gonic/gin"
	"github.com/Dungsenpai-ux/Practice_Go/service"
	"github.com/Dungsenpai-ux/Practice_Go/controller"
	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/Dungsenpai-ux/Practice_Go/config"
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
	m, err := migrate.New("file://migrations", "postgres://postgres:123456@localhost:5432/movies_db?sslmode=disable")
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

		movie := model.Movie{
			Title:  title,
			Year:   year,
			Genres: strings.Replace(record[2], "|", ", ", -1),
		}
		_, err = service.InsertMovie(ctx, movie)
		if err != nil {
			log.Printf("Bỏ qua: %v", err)
		}
	}
	return nil
}

func main() {
	if err := config.Connect(); err != nil {
		log.Fatal(err)
	}
	defer config.Close()

	if err := runMigrations(); err != nil {
		log.Fatal(err)
	}

	if err := seedData(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("POST /movies", controller.CreateMovie)
	http.HandleFunc("GET /movies/", controller.GetMovie)
	http.HandleFunc("GET /movies/search", controller.SearchMovies)

	log.Println("Máy chủ khởi động tại :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
// func main() {
// 	// Khởi tạo router Gin
// 	r := gin.Default()

// 	// Khởi tạo cấu hình
// 	cfg := config.NewConfig()

// 	// Đăng ký endpoint health
// 	healthController := controller.NewHealthController(cfg)
// 	r.GET("/healthz", healthController.HealthCheck)

// 	// Khởi chạy server
// 	r.Run(":8080")
// }