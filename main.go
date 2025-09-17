package main

import (
	// "log"
	"fmt"
	// "net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/controller"
	// "github.com/Dungsenpai-ux/Practice_Go/service/db/seed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations() error {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Xây dựng chuỗi DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// func main() {
// 	if err := config.Connect(); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer config.Close()

// 	if err := runMigrations(); err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := seedpkg.SeedData(); err != nil {
// 		log.Fatal(err)
// 	}

// 	http.HandleFunc("POST /movies", controller.CreateMovie)
// 	http.HandleFunc("GET /movies/", controller.GetMovie)
// 	http.HandleFunc("GET /movies/search", controller.SearchMovies)

// 	log.Println("Máy chủ khởi động tại :8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {
	// Khởi tạo router Gin
	r := gin.Default()

	// Khởi tạo cấu hình
	cfg := config.NewConfig()

	// Đăng ký endpoint health
	healthController := controller.NewHealthController(cfg)
	r.GET("/healthz", healthController.HealthCheck)

	// Khởi chạy server
	r.Run(":8080")
}
