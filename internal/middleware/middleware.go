package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			log.Println("Unauthorized access attempt")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		log.Printf("Authorized request with token: %s\n", token)
		c.Next()
	}
}
