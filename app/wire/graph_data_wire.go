//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/graph_data_handlers"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitializeCashDepositsHandler(postgresClient *gorm.DB, redisClient *redis.Client) *graph_data_handlers.CashDeposits {
	wire.Build(
		// handler
		graph_data_handlers.NewCashDeposits,
		// repositories
		repositories.NewTransactionRepository,
	)

	return new(graph_data_handlers.CashDeposits)
}

func InitializeTotalHoldingsPerDayHandler(postgresClient *gorm.DB, redisClient *redis.Client) *graph_data_handlers.TotalHoldingsPerDay {
	wire.Build(
		// handler
		graph_data_handlers.NewTotalHoldingsPerDay,
		// repositories
		repositories.NewTransactionRepository,
		repositories.NewHistoricalDataRepository,
		repositories.NewAllocationRepository,
	)

	return new(graph_data_handlers.TotalHoldingsPerDay)
}

func InitializePercentageAllocationsHandler(postgresClient *gorm.DB, redisClient *redis.Client) *graph_data_handlers.PercentageAllocations {
	wire.Build(
		// handler
		graph_data_handlers.NewPercentageAllocations,
		// repositories
		repositories.NewAllocationRepository,
	)

	return new(graph_data_handlers.PercentageAllocations)
}

func InitializeTradesHandler(postgresClient *gorm.DB, redisClient *redis.Client) *graph_data_handlers.Trades {
	wire.Build(
		// handler
		graph_data_handlers.NewTrades,
		// repositories
		repositories.NewHistoricalDataRepository,
		repositories.NewTransactionRepository,
	)

	return new(graph_data_handlers.Trades)
}
