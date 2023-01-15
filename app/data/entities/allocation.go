package entities

type Allocation struct {
	Symbol     string  `bson:"symbol,omitempty"`
	Percentage float64 `bson:"percentage,omitempty"`
}

type Allocations struct {
	PortfolioId string       `bson:"_id,omitempty"`
	Entries     []Allocation `bson:"entries,omitempty"`
}
