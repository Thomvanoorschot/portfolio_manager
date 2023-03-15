package graph_data_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/time_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	start := time_utils.TruncateToDay(firstTransaction.TransactedAt)
	end := time_utils.TruncateToDay(time.Now())
	uniqueSymbols := handler.transactionRepository.GetUniqueSymbolsForPortfolio(uuid.MustParse(portfolioId))
	historicalData := handler.historicalDataRepository.GetBySymbols(uniqueSymbols)

	var resp [][]float64
	holdings := handler.getHoldingsPerDay(portfolioId, start, end)
	previousHistoricalData := map[string][]decimal.Decimal{}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		var dayPrice decimal.Decimal
		for symbol, holdingAmount := range holdings[d] {
			if holdingAmount.IsZero() {
				continue
			}
			adjustedClose := historicalData[symbol][d].AdjustedClose
			dayPrice = dayPrice.Add(holdingAmount)
			holdingPrice := holdingAmount.Mul(adjustedClose)
			dayPrice = dayPrice.Add(holdingPrice)
			previousHistoricalData[symbol] = append(previousHistoricalData[symbol], adjustedClose)
		}
		resp = append(resp, []float64{float64(d.UnixMilli()), dayPrice.InexactFloat64()})
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
	holdings map[time.Time]map[string]decimal.Decimal,
	end time.Time,
	previousHistoricalData map[string][]decimal.Decimal,
) {
	// TODO Fix unneeded fromfloat
	endingDaySum := decimal.NewFromFloat(resp[len(resp)-1][1])
	allocations := &entities.Allocations{
		PortfolioId: portfolioId,
		Total:       endingDaySum,
	}
	for symbol, holdingAmount := range holdings[end] {
		if holdingAmount.IsZero() {
			continue
		}
		endingDaySymbolTotalValue := previousHistoricalData[symbol][len(previousHistoricalData[symbol])-1]
		total := endingDaySymbolTotalValue.Mul(holdingAmount)
		allocations.Entries = append(allocations.Entries, &entities.AllocationEntry{
			Symbol:            symbol,
			PercentageOfTotal: total.Div(endingDaySum).Mul(decimal.NewFromInt(100)),
			Total:             total,
			Amount:            holdingAmount,
		})
	}
	handler.allocationRepository.Upsert(portfolioId, allocations)
}

func (handler *TotalHoldingsPerDay) processTransactions(transactions entities.Transactions) map[time.Time]entities.Transactions {
	mappedTransactions := map[time.Time]entities.Transactions{}
	for _, transaction := range transactions {
		truncatedTransactedAt := time_utils.TruncateToDay(transaction.TransactedAt)
		mappedTransactions[truncatedTransactedAt] = append(mappedTransactions[truncatedTransactedAt], transaction)
	}
	return mappedTransactions
}

func (handler *TotalHoldingsPerDay) getHoldingsPerDay(portfolioId string,
	start time.Time,
	end time.Time) map[time.Time]map[string]decimal.Decimal {
	transactions := handler.transactionRepository.GetHoldingsTransactions(uuid.MustParse(portfolioId))

	mappedTransactions := handler.processTransactions(transactions)
	holdings := map[time.Time]map[string]decimal.Decimal{}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		holdings[d] = map[string]decimal.Decimal{}
		transactions, _ := mappedTransactions[d]
		if d != start {
			copyOfPreviousDay := map[string]decimal.Decimal{}
			previousDayHoldings := holdings[d.AddDate(0, 0, -1)]
			for k, v := range previousDayHoldings {
				copyOfPreviousDay[k] = v
			}
			holdings[d] = copyOfPreviousDay
		}
		for _, transaction := range transactions {
			holdings[d][transaction.Symbol] = holdings[d][transaction.Symbol].Add(transaction.Amount)
		}
	}
	return holdings
}
