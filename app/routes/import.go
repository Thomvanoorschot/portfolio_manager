package routes

import (
	"github.com/Thomvanoorschot/portfolioManager/app/wire"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupImportRoutes(routerGroup *gin.RouterGroup,
	postgresClient *gorm.DB,
	redisClient *redis.Client,
) *gin.RouterGroup {

	degiroImport := wire.InitializeDegiroImportHandler(postgresClient)
	historicalDataImport := wire.InitializeHistoricalDataImportHandler(postgresClient, redisClient)
	r := routerGroup.Group("/import")
	{
		r.POST("/degiro", degiroImport.Handle)
		r.POST("/historical", historicalDataImport.Handle)
	}
	return r
}
