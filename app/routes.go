package main

import (
	"github.com/Thomvanoorschot/portfolioManager/app/charting"
	"github.com/Thomvanoorschot/portfolioManager/app/importing"
	"github.com/Thomvanoorschot/portfolioManager/app/infrastructure"
	"net/http"
)

func SetupRoutes(server *infrastructure.Server) {
	server.Router.HandleFunc("/degiro-import", func(_ http.ResponseWriter, request *http.Request) { importing.DegiroImportHandler(server, request) }).Methods(http.MethodPost)
	server.Router.HandleFunc("/historical-import", func(_ http.ResponseWriter, request *http.Request) {
		importing.HistoricalDataImportHandler(server, request)
	}).Methods(http.MethodPost)
	server.Router.HandleFunc("/deposits", func(response http.ResponseWriter, request *http.Request) {
		charting.CashDepositsHandler(server, request, response)
	}).Methods(http.MethodGet)
	server.Router.HandleFunc("/holdings", func(response http.ResponseWriter, request *http.Request) {
		charting.HistoricalDataHandler(server, request, response)
	}).Methods(http.MethodGet)
}
