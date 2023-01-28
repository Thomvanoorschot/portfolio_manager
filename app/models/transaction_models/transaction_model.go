package transaction_models

import (
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/google/uuid"
	"time"
)

type Model struct {
	Id                uuid.UUID             `json:"id"`
	TransactedAt      time.Time             `json:"transactedAt"`
	CurrencyCode      string                `json:"currencyCode"`
	TransactionType   enums.TransactionType `json:"transactionType"`
	Product           string                `json:"product"`
	Amount            float64               `json:"amount"`
	PriceInCents      int64                 `json:"priceInCents"`
	CommissionInCents int64                 `json:"commissionInCents"`
	Symbol            string                `json:"symbol"`
}
