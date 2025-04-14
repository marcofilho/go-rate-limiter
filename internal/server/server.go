package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcofilho/go-rate-limiter/internal/limiter"
	"github.com/marcofilho/go-rate-limiter/internal/middleware"
)

func StartHTTPServer(port string, rateLimiter *limiter.RateLimiter) error {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		var key string
		var isToken bool
		var tokenRequestLimit int

		token := c.GetHeader("API_KEY")
		if token != "" {
			key = token
			isToken = true
			tokenRequestLimit = rateLimiter.TokenRequestLimit
		} else {
			key = c.ClientIP()
			isToken = false
			tokenRequestLimit = 10
		}

		allowed, err := rateLimiter.Allow(key, isToken, tokenRequestLimit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			return
		}

		c.Next()
	})

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Request ACCEPTED"})
	})

	address := fmt.Sprintf(":%s", port)
	return router.Run(address)
}
