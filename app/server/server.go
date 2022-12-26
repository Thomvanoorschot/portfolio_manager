package server

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Webserver struct {
	UnitOfWork *repositories.UnitOfWork
	fasthttp.Server
}

func Create() *Webserver {
	dsn := "host=localhost user=postgres password=Welkom01! dbname=portfoliomanager sslmode=disable"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	_ = db.AutoMigrate(&entities.Portfolio{}, &entities.Transaction{}, &entities.HistoricalData{})
	unitOfWork := &repositories.UnitOfWork{
		TransactionRepository:    &repositories.TransactionRepository{DB: db},
		PortfolioRepository:      &repositories.PortfolioRepository{DB: db},
		HistoricalDataRepository: &repositories.HistoricalDataRepository{DB: db},
	}

	return &Webserver{UnitOfWork: unitOfWork}
}
