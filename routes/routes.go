package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/handlers"
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/holdings"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *server.Webserver) {
	server.Use(CORSMiddleware())
	server.POST("/degiro-import", func(ctx *gin.Context) {
		handlers.DegiroImportHandler(server, ctx)
	})
	server.POST("/historical-import", func(ctx *gin.Context) {
		handlers.HistoricalDataImportHandler(server, ctx)
	})
	server.GET("/deposits/:portfolioId", func(ctx *gin.Context) {
		handlers.CashDepositsHandler(server, ctx)
	})
	server.GET("/holdings/per-day/:portfolioId", func(ctx *gin.Context) {
		holdings.PerDayHandler(server, ctx)
	})
	server.GET("/holdings/allocation/:portfolioId", func(ctx *gin.Context) {
		holdings.AllocationHandler(server, ctx)
	})
	server.GET("/trades/:portfolioId", func(ctx *gin.Context) {
		handlers.TradesHandler(server, ctx)
	})
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
