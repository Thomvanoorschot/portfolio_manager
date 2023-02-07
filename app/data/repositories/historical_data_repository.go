package repositories

import (
	"context"
	"encoding/json"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/redis/go-redis/v9"
	"time"
)

type HistoricalDataRepository struct {
	Rdb *redis.Client
}

func ProvideHistoricalDataRepository(rdb *redis.Client) *HistoricalDataRepository {
	return &HistoricalDataRepository{rdb}
}

func (p *HistoricalDataRepository) GetBySymbols(symbols []string) map[string]map[time.Time]entities.HistoricalDataEntry {
	historicalDataList := map[string]map[time.Time]entities.HistoricalDataEntry{}
	result, _ := p.Rdb.HMGet(context.Background(), "historical", symbols...).Result()
	c := make(chan entities.HistoricalData)
	for _, res := range result {
		go func(res interface{}) {
			historicalData := entities.HistoricalData{}
			_ = json.Unmarshal([]byte(res.(string)), &historicalData)
			c <- historicalData
		}(res)
	}
	for i := 0; i < len(result); i++ {
		res := <-c
		historicalDataList[res.Symbol] = res.Entries
	}
	return historicalDataList
}

func (p *HistoricalDataRepository) GetLastBySymbol(symbols []string) map[string]*entities.HistoricalDataEntry {
	m := map[string]*entities.HistoricalDataEntry{}
	bySymbols := p.GetBySymbols(symbols)
	for s, dataPerSymbol := range bySymbols {
		for i := 0; i < 30; i++ {
			last := dataPerSymbol[helpers.TruncateToDay(time.Now().AddDate(0, 0, -i))]
			if &last != nil {
				m[s] = &last
				break
			}
		}
	}

	return m
}
func (p *HistoricalDataRepository) GetBySymbol(symbol string) *entities.HistoricalData {
	historicalData := &entities.HistoricalData{}
	_ = p.Rdb.HGet(context.Background(), "historical", symbol).Scan(historicalData)
	return nil
}
func (p *HistoricalDataRepository) Upsert(historicalData *entities.HistoricalData) {
	bytes, _ := json.Marshal(historicalData)
	_, err := p.Rdb.HSet(context.Background(), "historical", historicalData.Symbol, bytes).Result()
	if err != nil {
		return
	}
}
