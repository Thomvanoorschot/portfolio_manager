package transaction_models

import (
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Id              uuid.UUID             `json:"id"`
	TransactedAt    time.Time             `json:"transactedAt"`
	CurrencyCode    string                `json:"currencyCode"`
	TransactionType enums.TransactionType `json:"transactionType"`
	Product         string                `json:"product"`
	Amount          decimal.Decimal       `json:"amount"`
	Price           decimal.Decimal       `json:"priceInCents"`
	Symbol          string                `json:"symbol"`
}
