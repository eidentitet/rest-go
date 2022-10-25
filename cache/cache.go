package cache

import (
	memCache "github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

var Cache *memCache.Cache
var NoExpiration = memCache.NoExpiration

func init() {
	once := sync.Once{}
	once.Do(func() {
		log.Println("Creating cache..")
		if Cache == nil {
			Cache = memCache.New(10*time.Minute, 5*time.Minute)
		}
	})
}

func Memory() *memCache.Cache {
	return Cache
}
