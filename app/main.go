package main

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/infrastructure"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	router := mux.NewRouter()

	dsn := "host=localhost user=postgres password=Welkom01! dbname=portfoliomanager sslmode=disable"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	_ = db.AutoMigrate(&entities.Portfolio{}, &entities.Transaction{}, &entities.HistoricalData{})
	unitOfWork := &data.UnitOfWork{
		TransactionRepository:    &repositories.TransactionRepository{DB: db},
		PortfolioRepository:      &repositories.PortfolioRepository{DB: db},
		HistoricalDataRepository: &repositories.HistoricalDataRepository{DB: db},
	}
	server := infrastructure.NewServer(unitOfWork, router)
	SetupRoutes(server)

	server.Run()
}
