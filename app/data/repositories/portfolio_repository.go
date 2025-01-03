package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioRepository struct {
	DB *gorm.DB
}

func NewPortfolioRepository(DB *gorm.DB) *PortfolioRepository {
	return &PortfolioRepository{DB: DB}
}

func (p *PortfolioRepository) Create(portfolio *entities.Portfolio) {
	p.DB.Create(portfolio)
}

func (p *PortfolioRepository) GetIncludingTransactions(portfolioId uuid.UUID) *entities.Portfolio {
	portfolio := &entities.Portfolio{}
	p.DB.Preload("Transactions").Where("id = ?", portfolioId).Find(portfolio)
	return portfolio
}
