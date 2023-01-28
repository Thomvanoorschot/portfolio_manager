package entities

import (
	"errors"
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/google/uuid"
	"time"
)

func ConvertToTransactionType(transactionType string) (enums.TransactionType, error) {
	switch transactionType {
	case "Koop":
		return enums.Purchase, nil
	case "Verkoop":
		return enums.Sale, nil
	default:
		return enums.TransactionTypeUnknown, errors.New("could not convert transactionType")
	}
}

type Transaction struct {
	EntityBase
	PortfolioID       uuid.UUID
	TransactedAt      time.Time
	Symbol            string
	ISIN              string
	Product           string
	CurrencyCode      string
	Description       string
	Amount            float64
	PriceInCents      int64
	CommissionInCents int64
	ExternalId        string
	TransactionType   enums.TransactionType `gorm:"default:Unknown"`
	AssetType         enums.AssetType       `gorm:"default:Unknown"`
	UniqueHash        string
}

type Transactions []*Transaction
