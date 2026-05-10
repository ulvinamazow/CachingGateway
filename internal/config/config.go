package config

import "flag"

type Config struct {
	Port          string
	Origin        string
	RedisAddr     string
	RedisPassword string
}

func ParseFlags() (Config, bool) {
	port := flag.String("port", "3000", "port to run proxy server")
	origin := flag.String("origin", "", "origin server url")
	redisAddr := flag.String("redis", "localhost:6379", "redis server address")
	redisPassword := flag.String("redis-password", "", "redis password")
	clearCache := flag.Bool("clear-cache", false, "clear redis cache")

	flag.Parse()

	cfg := Config{
		Port:          *port,
		Origin:        *origin,
		RedisAddr:     *redisAddr,
		RedisPassword: *redisPassword,
	}
	return cfg, *clearCache
}
