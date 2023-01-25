package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/graph_data_handlers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
)

func GetGraphDataRoutes(routerGroup *gin.RouterGroup, server *server.Webserver) *gin.RouterGroup {
	r := routerGroup.Group("/graph")
	{
		r.GET("/deposits/:portfolioId", func(ctx *gin.Context) {
			graph_data_handlers.CashDepositsHandler(server, ctx)
		})
		r.GET("/holdings/per-day/:portfolioId", func(ctx *gin.Context) {
			graph_data_handlers.PerDayHandler(server, ctx)
		})
		r.GET("/holdings/allocation/:portfolioId", func(ctx *gin.Context) {
			graph_data_handlers.AllocationHandler(server, ctx)
		})
		r.GET("/trades/:portfolioId", func(ctx *gin.Context) {
			graph_data_handlers.TradesHandler(server, ctx)
		})
	}
	return r
}
