package server

import (
	"context"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	nosqlDb := client.Database("historicalData")

	unitOfWork := &repositories.UnitOfWork{
		TransactionRepository:    &repositories.TransactionRepository{DB: db},
		PortfolioRepository:      &repositories.PortfolioRepository{DB: db},
		HistoricalDataRepository: repositories.ProvideHistoricalDataRepository(nosqlDb),
	}

	return &Webserver{
		UnitOfWork: unitOfWork,
		Engine:     g,
	}
}
