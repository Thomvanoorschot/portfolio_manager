package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/wire"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionRoutes(routerGroup *gin.RouterGroup,
	postgresClient *gorm.DB,
) *gin.RouterGroup {
	getByPortfolioId := wire.InitializeGetByPortfolioIdHandler(postgresClient)
	update := wire.InitializeUpdateHandler(postgresClient)
	updateTransactionSymbols := wire.InitializeUpdateTransactionSymbolsHandler(postgresClient)
	r := routerGroup.Group("/transaction")
	{
		r.GET("/:portfolioId", getByPortfolioId.Handle)
		r.PUT("/update", update.Handle)
		r.PUT("/update-symbols", updateTransactionSymbols.Handle)
	}
	return r
}
