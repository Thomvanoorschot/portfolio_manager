package import_handlers

import (
	"encoding/csv"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/Thomvanoorschot/portfolioManager/app/time_utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gitlab.com/metakeule/fmtdate"
	"io"
	"log"
	"net/http"
	"time"
)

type HistoricalDataImport struct {
	transactionRepository    *repositories.TransactionRepository
	historicalDataRepository *repositories.HistoricalDataRepository
}

func NewHistoricalDataImport(transactionRepository *repositories.TransactionRepository,
	historicalDataRepository *repositories.HistoricalDataRepository) *HistoricalDataImport {
	return &HistoricalDataImport{transactionRepository: transactionRepository,
		historicalDataRepository: historicalDataRepository}
}

var currencyMap = map[string]string{
	"USD": "EUR=X",
	"EUR": "EURUSD=X",
}

func (handler *HistoricalDataImport) Handle(_ *gin.Context) {
	symbolAssetTypePairs := handler.transactionRepository.GetUniqueSymbolAssetTypePairs()
	for _, symbolAssetTypePair := range symbolAssetTypePairs {
		overriddenSymbol := symbolAssetTypePair.Symbol
		if symbolAssetTypePair.AssetType == enums.Cash {
			overriddenSymbol = symbolAssetTypePair.Symbol + "=X"
		}
		url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=1d&events=history&includeAdjustedClose=true",
			overriddenSymbol,
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
			Symbol:  symbolAssetTypePair.Symbol,
			Entries: map[time.Time]entities.HistoricalDataEntry{},
		}
		var lastConvertedTime time.Time
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
			timestamp, _ := fmtdate.Parse("YYYY-MM-DD", line[0])
			timestamp = timestamp.UTC()
			for d := timestamp; d.After(lastConvertedTime); d = d.AddDate(0, 0, -1) {
				if lastConvertedTime.IsZero() {
					break
				}
				if d.Equal(timestamp) {
					continue
				}
				convertLine(d, line, historicalData.Entries)
			}
			couldConvert := convertLine(timestamp, line, historicalData.Entries)
			if couldConvert {
				lastConvertedTime = timestamp
			}
		}
		now := time_utils.TruncateToDay(time.Now())
		for d := now; d.After(lastConvertedTime); d = d.AddDate(0, 0, -1) {
			if lastConvertedTime.IsZero() {
				break
			}
			historicalData.Entries[d] = historicalData.Entries[lastConvertedTime]
		}
		handler.historicalDataRepository.Upsert(&historicalData)
	}
}

func convertLine(timestamp time.Time, line []string, historicalDataList map[time.Time]entities.HistoricalDataEntry) bool {
	historicalData := entities.HistoricalDataEntry{}
	historicalData.Timestamp = timestamp
	open := decimal.RequireFromString(line[1])
	if open.IsZero() {
		return false
	}
	historicalData.Open = open
	historicalData.High = decimal.RequireFromString(line[2])
	historicalData.Low = decimal.RequireFromString(line[3])
	historicalData.Close = decimal.RequireFromString(line[4])
	historicalData.AdjustedClose = decimal.RequireFromString(line[5])
	historicalData.Volume = decimal.RequireFromString(line[6])
	historicalDataList[timestamp] = historicalData
	return true
}
