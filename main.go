package main

import (
	"log"

	"net/http"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/controller"
	"github.com/Dungsenpai-ux/Practice_Go/service/db"

	//"github.com/gin-gonic/gin"
)

	"github.com/gin-gonic/gin"
)

// func main() {
// 	if err := config.Connect(); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer config.Close()

// 	if err := db.SeedData(); err != nil {
// 		log.Fatal(err)
// 	}

// 	controller.RegisterHTTPRoutes()

// 	log.Println("Máy chủ khởi động tại :8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {
	if err := config.Connect(); err != nil {
		log.Fatal(err)
	}
	defer config.Close()

	if err := config.ConnectMemcached(); err != nil {
		log.Fatal(err)
	}
	defer config.CloseMemcached()



	if err := db.SeedData(); err != nil {
		log.Fatal(err)
	}

	controller.RegisterHTTPRoutes()
	log.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
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
	r := gin.Default()

	// Khởi tạo cấu hình và thiết lập routes tập trung
	cfg := config.NewConfig()
	controller.SetupRoutes(r, cfg)

	// Khởi chạy server
	r.Run(":8080")
}
