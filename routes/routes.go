package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/handlers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func SetupRoutes(server *server.Webserver) {
	r := router.New()
	r.POST("/degiro-import", func(ctx *fasthttp.RequestCtx) {
		handlers.DegiroImportHandler(server, ctx)
	})
	r.POST("/historical-import", func(ctx *fasthttp.RequestCtx) {
		handlers.HistoricalDataImportHandler(server, ctx)
	})
	r.GET("/deposits", func(ctx *fasthttp.RequestCtx) {
		handlers.CashDepositsHandler(server, ctx)
	})
	r.GET("/holdings", func(ctx *fasthttp.RequestCtx) {
		handlers.HistoricalDataHandler(server, ctx)
	})
	server.Handler = r.Handler
}
