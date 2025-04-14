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

		token := c.GetHeader("API_KEY")
		if token != "" {
			key = token
			isToken = true
		} else {
			key = c.ClientIP()
			isToken = false
		}

		allowed, err := rateLimiter.Allow(key, isToken)
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

	// Apply the authentication middleware to protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/secure-ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "secure pong"})
	})

	// Public route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	address := fmt.Sprintf(":%s", port)
	return router.Run(address)
}
