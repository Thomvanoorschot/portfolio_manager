package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"gorm.io/gorm"
)

type PortfolioRepository struct {
	DB *gorm.DB
}

func ProvidePortfolioRepository(DB *gorm.DB) PortfolioRepository {
	return PortfolioRepository{DB: DB}
}

func (p *PortfolioRepository) Create(portfolio *entities.Portfolio) {
	p.DB.Create(portfolio)
}
