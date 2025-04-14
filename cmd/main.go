package main

import (
	"log"
	"time"

	"github.com/marcofilho/go-rate-limiter/configs"
	"github.com/marcofilho/go-rate-limiter/internal/limiter"
	"github.com/marcofilho/go-rate-limiter/internal/server"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	storage := limiter.NewRedisStorage(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	rateLimiter := limiter.NewRateLimiter(
		storage,
		cfg.RateLimiterMaxRequests,
		time.Duration(cfg.RateLimiterBlockDuration)*time.Second,
		cfg.TokenRequestLimit,
	)

	port := cfg.WebServerPort
	if port == "" {
		port = "8080"
	}

	log.Fatal(server.StartHTTPServer(port, rateLimiter))
}
