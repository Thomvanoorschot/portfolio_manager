package entities

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Allocations struct {
	PortfolioId string             `json:"portfolioId"`
	Total       decimal.Decimal    `json:"total"`
	Entries     []*AllocationEntry `json:"entries"`
}

type AllocationEntry struct {
	Symbol            string          `json:"symbol"`
	PercentageOfTotal decimal.Decimal `json:"percentageOfTotal"`
	Total             decimal.Decimal `json:"total"`
	Amount            decimal.Decimal `json:"amount"`
}

func (m *Allocations) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
