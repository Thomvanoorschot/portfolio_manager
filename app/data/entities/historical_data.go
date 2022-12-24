package entities

import "time"

type HistoricalData struct {
	Ticker        string
	Timestamp     time.Time
	Open          float64
	High          float64
	Low           float64
	Close         float64
	AdjustedClose float64
	Volume        int
	EntityBase
}
