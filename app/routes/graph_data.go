package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/wire"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupGraphDataRoutes(routerGroup *gin.RouterGroup,
	postgresClient *gorm.DB,
	redisClient *redis.Client,
) *gin.RouterGroup {
	cashDeposits := wire.InitializeCashDepositsHandler(postgresClient, redisClient)
	totalHoldingsPerDay := wire.InitializeTotalHoldingsPerDayHandler(postgresClient, redisClient)
	percentageAllocations := wire.InitializePercentageAllocationsHandler(postgresClient, redisClient)
	trades := wire.InitializeTradesHandler(postgresClient, redisClient)
	r := routerGroup.Group("/graph")
	{
		r.GET("/deposits/:portfolioId", cashDeposits.Handle)
		r.GET("/holdings/total/per-day/:portfolioId", totalHoldingsPerDay.Handle)
		//r.GET("/-holdings/symbol/:symbol/per-day/:portfolioId", func(ctx *gin.Context) {
		//	graph_data_handlers.HoldingsForSymbolPerDay(server, ctx)
		//})
		r.GET("/holdings/allocation/:portfolioId", percentageAllocations.Handle)
		r.GET("/trades/:portfolioId", trades.Handle)
	}
	return r
}
