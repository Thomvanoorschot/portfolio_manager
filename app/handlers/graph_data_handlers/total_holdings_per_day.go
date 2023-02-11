package graph_data_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"time"
)

type TotalHoldingsPerDay struct {
	transactionRepository    *repositories.TransactionRepository
	historicalDataRepository *repositories.HistoricalDataRepository
	allocationRepository     *repositories.AllocationRepository
}

func NewTotalHoldingsPerDay(transactionRepository *repositories.TransactionRepository,
	historicalDataRepository *repositories.HistoricalDataRepository,
	allocationRepository *repositories.AllocationRepository,
) *TotalHoldingsPerDay {
	return &TotalHoldingsPerDay{
		transactionRepository:    transactionRepository,
		historicalDataRepository: historicalDataRepository,
		allocationRepository:     allocationRepository,
	}
}

func (handler *TotalHoldingsPerDay) Handle(ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	transactions := handler.transactionRepository.GetHoldingsTransactions(uuid.MustParse(portfolioId))
	if len(transactions) == 0 {
		return
	}

	firstTransaction := transactions[0]
	start := helpers.TruncateToDay(firstTransaction.TransactedAt)
	end := helpers.TruncateToDay(time.Now())
	uniqueSymbols := handler.transactionRepository.GetUniqueSymbolsForPortfolio(uuid.MustParse(portfolioId))
	historicalData := handler.historicalDataRepository.GetBySymbols(uniqueSymbols)

	var resp [][]float64
	holdings := handler.getHoldingsPerDay(portfolioId, start, end)
	previousHistoricalData := map[string][]float64{}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		var dayPrice float64
		for symbol, h := range holdings[d] {
			adjustedClose := historicalData[symbol][d].AdjustedClose
			if adjustedClose != 0 {
				previousHistoricalData[symbol] = append(previousHistoricalData[symbol], adjustedClose)
			} else {
				adjustedClose = previousHistoricalData[symbol][len(previousHistoricalData[symbol])-1]
			}
			dayPrice += h * adjustedClose
		}
		resp = append(resp, []float64{float64(d.UnixMilli()), math.Round(dayPrice*100) / 100})
	}

	handler.persistAllocations(resp,
		portfolioId,
		holdings,
		end,
		previousHistoricalData)

	ctx.JSON(http.StatusOK, resp)
}

func (handler *TotalHoldingsPerDay) persistAllocations(resp [][]float64,
	portfolioId string,
	holdings map[time.Time]map[string]float64,
	end time.Time,
	previousHistoricalData map[string][]float64,
) {
	endingDaySum := resp[len(resp)-1][1]
	allocations := &entities.Allocations{
		PortfolioId: portfolioId,
		Total:       endingDaySum,
	}
	for symbol, h := range holdings[end] {
		if h == 0 {
			continue
		}
		endingDaySymbolTotalValue := previousHistoricalData[symbol][len(previousHistoricalData[symbol])-1]
		total := endingDaySymbolTotalValue * h
		allocations.Entries = append(allocations.Entries, &entities.AllocationEntry{
			Symbol:     symbol,
			Percentage: total / endingDaySum * 100,
			Total:      total,
			Amount:     h,
		})
	}
	handler.allocationRepository.Upsert(portfolioId, allocations)
}

func (handler *TotalHoldingsPerDay) processTransactions(transactions entities.Transactions) map[time.Time]entities.Transactions {
	mappedTransactions := map[time.Time]entities.Transactions{}
	for _, transaction := range transactions {
		truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)
		mappedTransactions[truncatedTransactedAt] = append(mappedTransactions[truncatedTransactedAt], transaction)
	}
	return mappedTransactions
}

func (handler *TotalHoldingsPerDay) getHoldingsPerDay(portfolioId string,
	start time.Time,
	end time.Time) map[time.Time]map[string]float64 {
	transactions := handler.transactionRepository.GetHoldingsTransactions(uuid.MustParse(portfolioId))

	mappedTransactions := handler.processTransactions(transactions)
	holdings := map[time.Time]map[string]float64{}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		holdings[d] = map[string]float64{}
		transactions, _ := mappedTransactions[d]
		if d != start {
			copyOfPreviousDay := map[string]float64{}
			previousDayHoldings := holdings[d.AddDate(0, 0, -1)]
			for k, v := range previousDayHoldings {
				copyOfPreviousDay[k] = v
			}
			holdings[d] = copyOfPreviousDay
		}
		for _, transaction := range transactions {
			holdings[d][transaction.Symbol] += transaction.Amount
		}
	}
	return holdings
}
