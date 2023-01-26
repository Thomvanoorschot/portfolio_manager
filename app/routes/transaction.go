package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/transaction_handlers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
)

func GetTransactionRoutes(routerGroup *gin.RouterGroup, server *server.Webserver) *gin.RouterGroup {
	r := routerGroup.Group("/transaction")
	{
		r.GET("/:portfolioId", func(ctx *gin.Context) {
			transaction_handlers.GetByPortfolioId(server, ctx)
		})
		r.PUT("/update", func(ctx *gin.Context) {
			transaction_handlers.Update(server, ctx)
		})
		r.PUT("/update-symbols", func(ctx *gin.Context) {
			transaction_handlers.UpdateTransactionSymbols(server, ctx)
		})
	}
	return r
}
