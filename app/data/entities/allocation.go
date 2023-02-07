package entities

import "encoding/json"

type Allocations struct {
	PortfolioId string             `json:"portfolioId"`
	Total       float64            `json:"total"`
	Entries     []*AllocationEntry `json:"entries"`
}

type AllocationEntry struct {
	Symbol     string  `json:"symbol"`
	Percentage float64 `json:"percentage"`
	Total      float64 `json:"total"`
	Amount     float64 `json:"amount"`
}

func (m *Allocations) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
