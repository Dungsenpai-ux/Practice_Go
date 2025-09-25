package main

import (
	"log"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/controller"
	srv "github.com/Dungsenpai-ux/Practice_Go/service"
	dbrepo "github.com/Dungsenpai-ux/Practice_Go/service/db"
	"github.com/Dungsenpai-ux/Practice_Go/service/external"
)

// Main: tối giản theo Option B – chỉ gọi các hàm init, route nằm trong controller.
func main() {
	cfg := config.Load()

	pool, err := dbrepo.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := dbrepo.SeedData(pool); err != nil {
		log.Fatal(err)
	}

	cache := external.InitMemcache(cfg.MemcachedAddr)
	movieService := srv.NewMovieService(pool)
	movieHandler := controller.NewMovieHandler(movieService, cache)
	healthHandler := controller.NewHealthController(cfg)

	r := controller.BuildRouter(cfg, healthHandler, movieHandler)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
