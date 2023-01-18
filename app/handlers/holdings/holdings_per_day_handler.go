package holdings

import (
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"sync"
	"time"
)

type holding struct {
	amount                 float64
	symbolPriceAtGivenTime float64
	total                  float64
}

func PerDayHandler(server *server.Webserver, ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	transactionRepository := server.UnitOfWork.TransactionRepository
	transactions := *transactionRepository.GetBuyAndSellTransactions(uuid.MustParse(portfolioId))
	if len(transactions) == 0 {
		return
	}

	firstTransaction := transactions[0]
	start := helpers.TruncateToDay(firstTransaction.TransactedAt)
	end := helpers.TruncateToDay(time.Now())
	uniqueSymbols := transactionRepository.GetUniqueSymbolsForPortfolio(uuid.MustParse(portfolioId))
	historicalDataPerSymbol := server.UnitOfWork.HistoricalDataRepository.GetBySymbols(uniqueSymbols)

	holdings := helpers.ThreadSafeMap[time.Time, helpers.ThreadSafeMap[string, holding]]{
		Entries: map[time.Time]*helpers.ThreadSafeMap[string, holding]{},
	}
	holdings.Entries[start] = &helpers.ThreadSafeMap[string, holding]{
		RWMutex: sync.RWMutex{},
		Entries: map[string]*holding{},
	}
	var resp [][]float64
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		for _, transaction := range transactions {
			truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)
			if !truncatedTransactedAt.Equal(d) {
				continue
			}
			thisDaysHoldings := holdings.Entries[truncatedTransactedAt]
			thisDaysSymbolHoldings := thisDaysHoldings.Entries[transaction.Symbol]
			if thisDaysSymbolHoldings == nil {
				newHolding := &holding{}
				thisDaysHoldings.Entries[transaction.Symbol] = newHolding
				thisDaysSymbolHoldings = newHolding
			}
			thisDaysSymbolHoldings.amount += transaction.Amount
		}

		var dayPrice float64
		newSafeHoldings := &helpers.ThreadSafeMap[string, holding]{
			RWMutex: sync.RWMutex{},
			Entries: map[string]*holding{},
		}
		wg := sync.WaitGroup{}
		currentHoldings := holdings.Entries[d].Entries
		c := make(chan float64, len(currentHoldings))
		for symbol, h := range currentHoldings {
			wg.Add(1)
			go func(symbol string, h *holding, c chan float64) {
				defer wg.Done()
				var symbolPriceAtGivenTime float64
				for _, historicalData := range historicalDataPerSymbol[symbol] {
					if historicalData.Timestamp.Year() == d.Year() && historicalData.Timestamp.Month() == d.Month() && historicalData.Timestamp.Day() == d.Day() {
						symbolPriceAtGivenTime = historicalData.Close
						break
					}
				}
				if symbolPriceAtGivenTime == 0 {
					previousDaysHolding := holdings.Get(d.AddDate(0, 0, -1))
					s := previousDaysHolding.Get(symbol)
					if s != nil {
						symbolPriceAtGivenTime = s.symbolPriceAtGivenTime
					} else {
						fmt.Println("Could not find price")
					}
					// TODO Deal with symbol changes not having historical data
				}
				h.symbolPriceAtGivenTime = symbolPriceAtGivenTime
				h.total = symbolPriceAtGivenTime * h.amount
				c <- h.total
				newSafeHoldings.Add(symbol, &holding{
					amount:                 h.amount,
					symbolPriceAtGivenTime: h.symbolPriceAtGivenTime,
					total:                  h.total,
				})
			}(symbol, h, c)
		}
		wg.Wait()
		close(c)
		for elem := range c {
			dayPrice += elem
		}
		holdings.Entries[d.AddDate(0, 0, 1)] = newSafeHoldings

		resp = append(resp, []float64{float64(d.UnixMilli()), math.Round(dayPrice*100) / 100})
	}
	allocations := entities.Allocations{}

	var amountSum float64
	for _, h := range holdings.Entries[end].Entries {
		amountSum += h.total
	}
	for symbol, h := range holdings.Entries[end].Entries {
		if h.amount == 0 {
			continue
		}
		allocations.Entries = append(allocations.Entries, entities.Allocation{
			Symbol:     symbol,
			Percentage: h.total / amountSum * 100,
			Total:      h.total,
		})
	}
	server.UnitOfWork.AllocationRepository.Upsert(portfolioId, allocations)

	ctx.JSON(http.StatusOK, resp)
}
