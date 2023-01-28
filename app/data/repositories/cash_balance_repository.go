package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CashBalanceRepository struct {
	DB *gorm.DB
}

func (p *CashBalanceRepository) Update(cashBalanceId uuid.UUID, currencyCode string, amountInCents int64) {
	p.DB.Model(&entities.CashBalance{}).Where("id = ? AND currency_code = ?", cashBalanceId).Update("amount_in_cents", amountInCents)
}
