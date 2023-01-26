package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/import_handlers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
)

func GetImportRoutes(routerGroup *gin.RouterGroup, server *server.Webserver) *gin.RouterGroup {
	r := routerGroup.Group("/import")
	{
		r.POST("/degiro", func(ctx *gin.Context) {
			import_handlers.DegiroImport(server, ctx)
		})
		r.POST("/historical", func(ctx *gin.Context) {
			import_handlers.HistoricalDataImport(server, ctx)
		})
	}
	return r
}
