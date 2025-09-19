package config

import (
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/rs/zerolog/log"
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
		log.Warn().Err(err).Str("addr", addr).Msg("Memcached not available, continuing without cache")
		Memcached = nil
		return nil
	}
	_ = Memcached.Delete("test_connection")
	return nil
}

func CloseMemcached() {
}
