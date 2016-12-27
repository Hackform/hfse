package cache

import (
	"github.com/Hackform/hfse/service"
	"gopkg.in/redis.v5"
)

type (
	CacheInfo struct {
		url, pass string
		name      int
	}

	Cache struct {
		service.ServiceBase
		CacheInfo
		database *redis.Client
	}
)

func New(url string, name int, pass string) *Cache {
	return &Cache{
		CacheInfo: CacheInfo{
			url:  url,
			name: name,
			pass: pass,
		},
		database: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: pass,
			DB:       name,
		}),
	}
}

func (c *Cache) Start() {
	done := c.Connect()
	<-done
}

func (c *Cache) Shutdown() {
	c.Close()
}

func (c *Cache) Connect() <-chan bool {
	done := make(chan bool)
	go func() {
		_, err := c.database.Ping().Result()
		done <- err == nil
	}()
	return done
}

func (c *Cache) Close() {
	c.database.Close()
}
