package main

import (
	"fmt"

	"github.com/ulvinamazow/caching-gateway/internal/cache"
	"github.com/ulvinamazow/caching-gateway/internal/config"
	"github.com/ulvinamazow/caching-gateway/internal/middleware/metrics"
	"github.com/ulvinamazow/caching-gateway/internal/proxy"
	"github.com/ulvinamazow/caching-gateway/internal/redis"
	"github.com/ulvinamazow/caching-gateway/internal/server"
)

func main() {
	cfg, clearCache := config.ParseFlags()

	client := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)

	metrics.RegisterCollectors()

	cacheStore := cache.NewRedisCache(client)

	if clearCache {
		err := cacheStore.Clear()
		if err != nil {
			panic(err)
		}
		fmt.Println("cache cleared successfully")
		return
	}

	handler := proxy.NewHandler(cacheStore, cfg.Origin)

	server := server.NewServer(cfg.Port, handler)

	fmt.Println("proxy server starting on port: ", cfg.Port)
	fmt.Println("origin server: ", cfg.Origin)

	err := server.Start()
	if err != nil {
		panic(err)
	}

}
