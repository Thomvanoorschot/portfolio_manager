package entities

type Allocation struct {
	Symbol     string  `bson:"symbol,omitempty"`
	Percentage float64 `bson:"percentage,omitempty"`
	Total      float64 `bson:"total,omitempty"`
}

type Allocations struct {
	PortfolioId string       `bson:"_id,omitempty"`
	Total       float64      `bson:"total"`
	Entries     []Allocation `bson:"entries,omitempty"`
}
