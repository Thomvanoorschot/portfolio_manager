package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupV1Routes(r *gin.Engine,
	postgresClient *gorm.DB,
	redisClient *redis.Client) {
	r.Use(CORSMiddleware())

	v1 := r.Group("/api/v1")
	SetupGraphDataRoutes(v1, postgresClient, redisClient)
	SetupImportRoutes(v1, postgresClient, redisClient)
	SetupTransactionRoutes(v1, postgresClient)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
