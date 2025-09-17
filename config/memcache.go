package config

import (
	"log"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

var Memcached *memcache.Client

func ConnectMemcached() error {
	addr := os.Getenv("MEMCACHED_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11211"
	}

	Memcached = memcache.New(addr)
	// Kiểm tra kết nối
	err := Memcached.Set(&memcache.Item{
		Key:   "test_connection",
		Value: []byte("test"),
	})
	if err != nil {
		// Nếu không kết nối được, tắt cache và tiếp tục không dùng cache
		log.Printf("Memcached not available at %s: %v (continuing without cache)", addr, err)
		Memcached = nil
		return nil
	}
	_ = Memcached.Delete("test_connection")
	return nil
}

func CloseMemcached() {
}
