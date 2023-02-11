//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/handlers/transaction_handlers"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeGetByPortfolioIdHandler(postgresClient *gorm.DB) *transaction_handlers.GetByPortfolioId {
	wire.Build(
		// handler
		transaction_handlers.NewGetByPortfolioId,
		// repositories
		repositories.NewTransactionRepository,
	)

	return new(transaction_handlers.GetByPortfolioId)
}
func InitializeUpdateHandler(postgresClient *gorm.DB) *transaction_handlers.Update {
	wire.Build(
		// handler
		transaction_handlers.NewUpdate,
		// repositories
		repositories.NewTransactionRepository,
	)

	return new(transaction_handlers.Update)
}
func InitializeUpdateTransactionSymbolsHandler(postgresClient *gorm.DB) *transaction_handlers.UpdateTransactionSymbols {
	wire.Build(
		// handler
		transaction_handlers.NewUpdateTransactionSymbols,
		// repositories
		repositories.NewTransactionRepository,
	)

	return new(transaction_handlers.UpdateTransactionSymbols)
}
