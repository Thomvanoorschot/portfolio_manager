package repositories

import (
	"context"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type HistoricalDataRepository struct {
	Collection *mongo.Collection
}

func ProvideHistoricalDataRepository(database *mongo.Database) *HistoricalDataRepository {
	historicalDataCollection := database.Collection("historicalData")
	return &HistoricalDataRepository{Collection: historicalDataCollection}
}

func (p *HistoricalDataRepository) GetBySymbols(symbols []string) map[string]map[time.Time]*entities.HistoricalDataEntry {
	var historicalData []entities.HistoricalData
	filter := bson.M{"_id": bson.M{"$in": symbols}}

	find, _ := p.Collection.Find(context.TODO(), filter)
	_ = find.All(context.TODO(), &historicalData)

	m := make(map[string]map[time.Time]*entities.HistoricalDataEntry)
	for _, d := range historicalData {
		m[d.Symbol] = d.Entries
	}
	return m
}

func (p *HistoricalDataRepository) GetLastBySymbol(symbols []string) *helpers.ThreadSafeMap[string, entities.HistoricalDataEntry] {
	var historicalData []entities.HistoricalData
	filter := bson.M{"_id": bson.M{"$in": symbols}}

	find, _ := p.Collection.Find(context.TODO(), filter)
	_ = find.All(context.TODO(), &historicalData)

	m := helpers.ThreadSafeMap[string, entities.HistoricalDataEntry]{}
	m.Entries = map[string]*entities.HistoricalDataEntry{}
	for _, d := range historicalData {
		if d.Entries == nil {
			continue
		}
		for i := 0; i < 30; i++ {
			last := d.Entries[helpers.TruncateToDay(time.Now().AddDate(0, 0, -i))]
			if last != nil {
				m.Entries[d.Symbol] = last
				break
			}
		}
	}
	return &m
}
func (p *HistoricalDataRepository) GetBySymbol(symbol string) *entities.HistoricalData {
	var historicalData entities.HistoricalData
	filter := bson.D{{"_id", symbol}}

	_ = p.Collection.FindOne(context.TODO(), filter).Decode(&historicalData)
	return &historicalData
}
func (p *HistoricalDataRepository) Insert(historicalData *entities.HistoricalData) {
	filter := bson.M{"_id": historicalData.Symbol}
	upsert := true
	_, err := p.Collection.ReplaceOne(context.TODO(), filter, historicalData, &options.ReplaceOptions{Upsert: &upsert})
	if err != nil {
		fmt.Println(err)
	}
}
