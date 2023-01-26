package import_handlers

import (
	"encoding/csv"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"gitlab.com/metakeule/fmtdate"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HistoricalDataImport(server *server.Webserver, _ *gin.Context) {
	symbols := server.UnitOfWork.TransactionRepository.GetUniqueSymbols()
	for _, symbol := range symbols {
		url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=1d&events=history&includeAdjustedClose=true",
			symbol,
			time.Now().AddDate(-10, 0, 0).Unix(),
			time.Now().Unix())

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Could not get historical data for ticker", "", err)
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		reader := csv.NewReader(resp.Body)
		firstLineProcessed := false

		historicalData := entities.HistoricalData{
			Symbol:  symbol,
			Entries: map[time.Time]*entities.HistoricalDataEntry{},
		}

		for {
			line, readError := reader.Read()
			if readError == io.EOF {
				break
			} else if readError != nil {
				log.Fatal(readError)
			} else if !firstLineProcessed {
				firstLineProcessed = true
				continue
			}
			convertLine(line, historicalData.Entries)
		}
		server.UnitOfWork.HistoricalDataRepository.Insert(&historicalData)
	}
}

func convertLine(line []string, historicalDataList map[time.Time]*entities.HistoricalDataEntry) {
	timestamp, _ := fmtdate.Parse("YYYY-MM-DD", line[0])
	timestamp = timestamp.UTC()
	historicalData := entities.HistoricalDataEntry{}
	historicalData.Timestamp = timestamp
	open, _ := strconv.ParseFloat(line[1], 64)
	historicalData.Open = open
	high, _ := strconv.ParseFloat(line[2], 64)
	historicalData.High = high
	low, _ := strconv.ParseFloat(line[3], 64)
	historicalData.Low = low
	closeAmount, _ := strconv.ParseFloat(line[4], 64)
	historicalData.Close = closeAmount
	adjustedClose, _ := strconv.ParseFloat(line[5], 64)
	historicalData.AdjustedClose = adjustedClose
	volume, _ := strconv.Atoi(line[6])
	historicalData.Volume = volume
	historicalDataList[timestamp] = &historicalData
}
