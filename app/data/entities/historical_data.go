package entities

import (
	"time"
)

type HistoricalData struct {
	Symbol  string                            `json:"symbol"`
	Entries map[time.Time]HistoricalDataEntry `json:"entries"`
}

type HistoricalDataEntry struct {
	Timestamp     time.Time `json:"timestamp"`
	Open          float64   `json:"open"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Close         float64   `json:"close"`
	AdjustedClose float64   `json:"adjustedClose"`
	Volume        int       `json:"volume"`
}
