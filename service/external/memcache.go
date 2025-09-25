package external

import (
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

// InitMemcache initializes a memcache client; returns nil (no error) if unreachable to allow graceful degradation.
func InitMemcache(addr string) *memcache.Client {
	c := memcache.New(addr)
	if err := c.Set(&memcache.Item{Key: "boot_probe", Value: []byte("ok"), Expiration: 5}); err != nil {
		log.Printf("Memcached not available at %s: %v (continuing without cache)", addr, err)
		return nil
	}
	_ = c.Delete("boot_probe")
	return c
}
