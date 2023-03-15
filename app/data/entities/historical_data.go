package entities

import (
	"github.com/shopspring/decimal"
	"time"
)

type HistoricalData struct {
	Symbol  string                            `json:"symbol"`
	Entries map[time.Time]HistoricalDataEntry `json:"entries"`
}

type HistoricalDataEntry struct {
	Timestamp     time.Time       `json:"timestamp"`
	Open          decimal.Decimal `json:"open"`
	High          decimal.Decimal `json:"high"`
	Low           decimal.Decimal `json:"low"`
	Close         decimal.Decimal `json:"close"`
	AdjustedClose decimal.Decimal `json:"adjustedClose"`
	Volume        decimal.Decimal `json:"volume"`
}
