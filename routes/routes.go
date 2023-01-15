package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/handlers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *server.Webserver) {
	server.POST("/degiro-import", func(ctx *gin.Context) {
		handlers.DegiroImportHandler(server, ctx)
	})
	server.POST("/historical-import", func(ctx *gin.Context) {
		handlers.HistoricalDataImportHandler(server, ctx)
	})
	server.GET("/deposits/:portfolioId", func(ctx *gin.Context) {
		handlers.CashDepositsHandler(server, ctx)
	})
	server.GET("/holdings/:portfolioId", func(ctx *gin.Context) {
		handlers.HoldingsHandler(server, ctx)
	})
	server.GET("/trades/:portfolioId", func(ctx *gin.Context) {
		handlers.TradesHandler(server, ctx)
	})
}
