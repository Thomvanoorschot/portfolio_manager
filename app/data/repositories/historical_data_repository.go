package repositories

import (
	"context"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type HistoricalDataRepository struct {
	Collection *mongo.Collection
}

func ProvideHistoricalDataRepository(database *mongo.Database) *HistoricalDataRepository {
	historicalDataCollection := database.Collection("historicalData")
	return &HistoricalDataRepository{Collection: historicalDataCollection}
}

func (p *HistoricalDataRepository) GetBySymbols(symbols []string) map[string][]entities.HistoricalDataEntry {
	m := map[string][]entities.HistoricalDataEntry{}
	wg := sync.WaitGroup{}
	c := make(chan *entities.HistoricalData, len(symbols))

	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string, c chan *entities.HistoricalData) {
			historicalData := p.GetBySymbol(symbol)
			c <- historicalData
			wg.Done()
		}(symbol, c)
	}
	wg.Wait()
	close(c)
	for data := range c {
		m[data.Symbol] = data.Entries
	}
	return m
}
func (p *HistoricalDataRepository) GetBySymbol(symbol string) *entities.HistoricalData {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var historicalData entities.HistoricalData
	filter := bson.D{{"_id", symbol}}

	_ = p.Collection.FindOne(ctx, filter).Decode(&historicalData)
	return &historicalData
}
func (p *HistoricalDataRepository) Insert(historicalData *entities.HistoricalData) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"symbol", historicalData.Symbol}}
	upsert := true
	_, err := p.Collection.ReplaceOne(ctx, filter, historicalData, &options.ReplaceOptions{Upsert: &upsert})
	fmt.Println(err)
}
