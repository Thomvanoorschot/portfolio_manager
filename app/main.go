package main

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/routes"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	g := gin.Default()

	dsn := "host=localhost user=postgres password=Welkom01! dbname=portfoliomanager sslmode=disable"
	postgresClient, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	_ = postgresClient.AutoMigrate(&entities.Portfolio{}, &entities.Transaction{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	routes.SetupV1Routes(g, postgresClient, redisClient)
	log.Fatal(g.RunTLS("127.0.0.1:8000", "localhost.crt", "localhost.key"))
}
