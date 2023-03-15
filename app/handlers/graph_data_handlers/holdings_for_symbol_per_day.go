package graph_data_handlers

//
//import (
//	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
//	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
//	"github.com/Thomvanoorschot/portfolioManager/app/server"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	"math"
//	"net/http"
//	"time"
//)
//
//func HoldingsForSymbolPerDay(server *server.Webserver, ctx *gin.Context) {
//	portfolioId := ctx.Param("portfolioId")
//	symbol := ctx.Param("symbol")
//
//	transactionRepository := server.UnitOfWork.TransactionRepository
//	transactions := transactionRepository.GetHoldingsTransactionsPerSymbol(uuid.MustParse(portfolioId), symbol)
//	if len(transactions) == 0 {
//		return
//	}
//
//	firstTransaction := transactions[0]
//	start := helpers.TruncateToDay(firstTransaction.TransactedAt)
//	end := helpers.TruncateToDay(time.Now())
//	historicalData := server.UnitOfWork.HistoricalDataRepository.GetBySymbol(symbol)
//
//	var resp [][]float64
//	holdings := getHoldingsPerDay(server, portfolioId, start, end)
//	previousHistoricalData := map[string][]float64{}
//	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
//		var dayPrice float64
//		for symbol, h := range holdings[d] {
//			adjustedClose := historicalData[symbol][d].AdjustedCloseInFractionalUnits
//			if adjustedClose != 0 {
//				previousHistoricalData[symbol] = append(previousHistoricalData[symbol], adjustedClose)
//			} else {
//				adjustedClose = previousHistoricalData[symbol][len(previousHistoricalData[symbol])-1]
//			}
//			dayPrice += h * adjustedClose
//		}
//		resp = append(resp, []float64{float64(d.UnixMilli()), math.Round(dayPrice*100) / 100})
//	}
//
//	persistAllocations(resp,
//		portfolioId,
//		holdings,
//		end,
//		previousHistoricalData,
//		server)
//
//	ctx.JSON(http.StatusOK, resp)
//}
//
//func persistAllocations(resp [][]float64,
//	portfolioId string,
//	holdings map[time.Time]map[string]float64,
//	end time.Time,
//	previousHistoricalData map[string][]float64,
//	server *server.Webserver,
//) {
//	endingDaySum := resp[len(resp)-1][1]
//	allocations := &entities.Allocations{
//		PortfolioId: portfolioId,
//		TotalInFractionalUnits:       endingDaySum,
//	}
//	for symbol, h := range holdings[end] {
//		if h == 0 {
//			continue
//		}
//		endingDaySymbolTotalValue := previousHistoricalData[symbol][len(previousHistoricalData[symbol])-1]
//		total := endingDaySymbolTotalValue * h
//		allocations.Entries = append(allocations.Entries, &entities.AllocationEntry{
//			Symbol:     symbol,
//			PercentageOfTotal: total / endingDaySum * 100,
//			TotalInFractionalUnits:      total,
//			AmountInFractionalUnits:     h,
//		})
//	}
//	server.UnitOfWork.AllocationRepository.Upsert(portfolioId, allocations)
//}
//
//func processTransactions(transactions entities.Transactions) map[time.Time]entities.Transactions {
//	mappedTransactions := map[time.Time]entities.Transactions{}
//	for _, transaction := range transactions {
//		truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)
//		mappedTransactions[truncatedTransactedAt] = append(mappedTransactions[truncatedTransactedAt], transaction)
//	}
//	return mappedTransactions
//}
//
//func getHoldingsPerDay(server *server.Webserver,
//	portfolioId string,
//	start time.Time,
//	end time.Time) map[time.Time]map[string]float64 {
//	transactions := server.UnitOfWork.TransactionRepository.GetHoldingsTransactions(uuid.MustParse(portfolioId))
//
//	mappedTransactions := processTransactions(transactions)
//	holdings := map[time.Time]map[string]float64{}
//	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
//		holdings[d] = map[string]float64{}
//		transactions, _ := mappedTransactions[d]
//		if d != start {
//			copyOfPreviousDay := map[string]float64{}
//			previousDayHoldings := holdings[d.AddDate(0, 0, -1)]
//			for k, v := range previousDayHoldings {
//				copyOfPreviousDay[k] = v
//			}
//			holdings[d] = copyOfPreviousDay
//		}
//		for _, transaction := range transactions {
//			holdings[d][transaction.Symbol] += transaction.AmountInFractionalUnits
//		}
//	}
//	return holdings
//}
