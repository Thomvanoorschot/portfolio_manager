package entities

import (
	"time"
)

type HistoricalData struct {
	Symbol  string                             `bson:"_id,omitempty"`
	Entries map[time.Time]*HistoricalDataEntry `bson:"entries,omitempty"`
}

type HistoricalDataEntry struct {
	Timestamp     time.Time `bson:"timestamp,omitempty"`
	Open          float64   `bson:"open,omitempty"`
	High          float64   `bson:"high,omitempty"`
	Low           float64   `bson:"low,omitempty"`
	Close         float64   `bson:"close,omitempty"`
	AdjustedClose float64   `bson:"adjustedClose,omitempty"`
	Volume        int       `bson:"volume,omitempty"`
}
