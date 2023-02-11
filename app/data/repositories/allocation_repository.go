package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/redis/go-redis/v9"
)

type AllocationRepository struct {
	Rdb *redis.Client
}

func NewAllocationRepository(rdb *redis.Client) *AllocationRepository {
	return &AllocationRepository{Rdb: rdb}
}

func (p *AllocationRepository) Upsert(portfolioId string, allocations *entities.Allocations) {
	bytes, _ := json.Marshal(allocations)
	_, err := p.Rdb.HSet(context.Background(), "allocation", portfolioId, bytes).Result()
	if err != nil {
		return
	}
}
func (p *AllocationRepository) GetByPortfolioId(portfolioId string) *entities.Allocations {
	allocations := &entities.Allocations{}
	err := p.Rdb.HGet(context.Background(), "allocation", portfolioId).Scan(allocations)
	fmt.Println(err)
	return allocations
}
