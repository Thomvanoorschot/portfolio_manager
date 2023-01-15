package main

import (
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/Thomvanoorschot/portfolioManager/routes"
	"log"
)

func main() {
	webServer := server.Create()
	routes.SetupRoutes(webServer)
	log.Fatal(webServer.RunTLS("127.0.0.1:8000", "localhost.crt", "localhost.key"))
}
