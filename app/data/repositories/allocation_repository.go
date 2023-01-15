package repositories

import (
	"context"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AllocationRepository struct {
	Collection *mongo.Collection
}

func ProvideAllocationRepository(database *mongo.Database) *AllocationRepository {
	allocationCollection := database.Collection("allocation")
	return &AllocationRepository{Collection: allocationCollection}
}

func (p *AllocationRepository) Upsert(portfolioId string, allocations entities.Allocations) {
	filter := bson.M{"_id": portfolioId}
	upsert := true
	_, err := p.Collection.ReplaceOne(context.TODO(), filter, allocations, &options.ReplaceOptions{Upsert: &upsert})
	if err != nil {
		fmt.Println(err)
	}
}
func (p *AllocationRepository) GetByPortfolioId(portfolioId string) *entities.Allocations {
	var allocations entities.Allocations
	filter := bson.D{{"_id", portfolioId}}

	_ = p.Collection.FindOne(context.TODO(), filter).Decode(&allocations)
	return &allocations
}
