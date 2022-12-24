package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"gorm.io/gorm"
)

type HistoricalDataRepository struct {
	DB *gorm.DB
}

func ProvideHistoricalDataRepository(DB *gorm.DB) HistoricalDataRepository {
	return HistoricalDataRepository{DB: DB}
}

func (p *HistoricalDataRepository) GetBySymbols(symbols []string) map[string][]entities.HistoricalData {
	var historicalData []entities.HistoricalData
	p.DB.Where("ticker IN ?", symbols).Find(&historicalData)
	m := map[string][]entities.HistoricalData{}
	for _, d := range historicalData {
		m[d.Ticker] = append(m[d.Ticker], d)
	}
	return m
}
func (p *HistoricalDataRepository) GetBySymbol(symbol string) []entities.HistoricalData {
	var historicalData []entities.HistoricalData
	p.DB.Where("ticker = ?", symbol).Find(&historicalData)
	return historicalData
}
func (p *HistoricalDataRepository) BatchInsert(historicalData *[]entities.HistoricalData) {
	p.DB.Create(historicalData)
}
