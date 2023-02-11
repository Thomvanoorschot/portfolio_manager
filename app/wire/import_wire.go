//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/import_handlers"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitializeDegiroImportHandler(postgresClient *gorm.DB) *import_handlers.DegiroImport {
	wire.Build(
		// handler
		import_handlers.NewDegiroImport,
		// repositories
		repositories.NewTransactionRepository,
		repositories.NewPortfolioRepository,
	)

	return new(import_handlers.DegiroImport)
}
func InitializeHistoricalDataImportHandler(postgresClient *gorm.DB, redisClient *redis.Client) *import_handlers.HistoricalDataImport {
	wire.Build(
		// handler
		import_handlers.NewHistoricalDataImport,
		// repositories
		repositories.NewTransactionRepository,
		repositories.NewHistoricalDataRepository,
	)

	return new(import_handlers.HistoricalDataImport)
}
