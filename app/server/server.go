package server

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Webserver struct {
	UnitOfWork *repositories.UnitOfWork
	*gin.Engine
}

func Create() *Webserver {
	g := gin.Default()
	dsn := "host=localhost user=postgres password=Welkom01! dbname=portfoliomanager sslmode=disable"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	_ = db.AutoMigrate(&entities.Portfolio{}, &entities.Transaction{})

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	unitOfWork := &repositories.UnitOfWork{
		TransactionRepository:    &repositories.TransactionRepository{DB: db},
		PortfolioRepository:      &repositories.PortfolioRepository{DB: db},
		HistoricalDataRepository: repositories.ProvideHistoricalDataRepository(rdb),
		AllocationRepository:     repositories.ProvideAllocationRepository(rdb),
	}

	return &Webserver{
		UnitOfWork: unitOfWork,
		Engine:     g,
	}
}
