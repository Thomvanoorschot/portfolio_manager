package data

import "github.com/Thomvanoorschot/portfolioManager/app/data/repositories"

type UnitOfWork struct {
	TransactionRepository    *repositories.TransactionRepository
	PortfolioRepository      *repositories.PortfolioRepository
	HistoricalDataRepository *repositories.HistoricalDataRepository
}
