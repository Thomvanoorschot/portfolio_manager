package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioRepository struct {
	DB *gorm.DB
}

func (p *PortfolioRepository) Create(portfolio *entities.Portfolio) {
	p.DB.Create(portfolio)
}

func (p *PortfolioRepository) GetIncludingTransactionsAndCashBalances(portfolioId uuid.UUID) *entities.Portfolio {
	portfolio := &entities.Portfolio{}
	p.DB.Preload("Transactions").Preload("CashBalances").Where("id = ?", portfolioId).Find(portfolio)
	return portfolio
}
