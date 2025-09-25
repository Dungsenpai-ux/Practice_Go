package main

import (
	"log"
	"net/http"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/controller"
	"github.com/Dungsenpai-ux/Practice_Go/service"
	dbrepo "github.com/Dungsenpai-ux/Practice_Go/service/db"
	"github.com/Dungsenpai-ux/Practice_Go/service/external"
)

func main() {
	// 1. Load config
	cfg := config.Load()

	// 2. Init DB
	pool, err := dbrepo.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// 3. Seed data if empty
	if err := dbrepo.SeedData(pool); err != nil {
		log.Fatal(err)
	}

	// 4. Init external services (memcache)
	cache := external.InitMemcache(cfg.MemcachedAddr)

	// 5. Construct services
	movieService := service.NewMovieService(pool)

	// 6. Handlers
	movieHandler := controller.NewMovieHandler(movieService, cache)
	healthHandler := controller.NewHealthController(cfg)

	// 7. Register HTTP routes (net/http style)
	controller.RegisterHTTPRoutes(movieHandler)
	http.Handle("/healthz", healthHandler)

	// 8. Start server
	log.Printf("HTTP server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	if err := config.Connect(); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer config.Close()
//
// 	if err := db.SeedData(); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Khởi tạo router Gin
// 	r := gin.Default()
//
// 	// Khởi tạo cấu hình và thiết lập routes tập trung
// 	cfg := config.NewConfig()
// 	controller.SetupRoutes(r, cfg)
//
// 	// Khởi chạy server
// 	r.Run(":8080")
// }
// Khởi tạo router Gin
// 	r := gin.Default()

// 	// Khởi tạo cấu hình và thiết lập routes tập trung
// 	cfg := config.NewConfig()
// 	controller.SetupRoutes(r, cfg)

// 	// Khởi chạy server
// 	r.Run(":8080")
// }
