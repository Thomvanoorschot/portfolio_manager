package graph_data_models

import "github.com/shopspring/decimal"

type Allocation struct {
	Name string          `json:"name"`
	Y    decimal.Decimal `json:"y"`
}
