package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/transaction_handlers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
)

func GetTransactionRoutes(routerGroup *gin.RouterGroup, server *server.Webserver) *gin.RouterGroup {
	r := routerGroup.Group("/transaction")
	{
		r.POST("/update-symbols", func(ctx *gin.Context) {
			transaction_handlers.UpdateTransactionSymbolsHandler(server, ctx)
		})
	}
	return r
}
