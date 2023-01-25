package graph_data_handlers

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
	amount                        float64
	priceOfSymbolPriceAtGivenTime float64
	total                         float64
}

type readOp struct {
	date   time.Time
	symbol string
	resp   chan *holding
}
type writeOp struct {
	date   time.Time
	symbol string
	val    *holding
	resp   chan bool
}

func PerDayHandler(server *server.Webserver, ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	transactionRepository := server.UnitOfWork.TransactionRepository
	transactions := transactionRepository.GetBuyAndSellTransactions(uuid.MustParse(portfolioId))
	if len(transactions) == 0 {
		return
	}

	firstTransaction := transactions[0]
	start := helpers.TruncateToDay(firstTransaction.TransactedAt)
	end := helpers.TruncateToDay(time.Now())
	uniqueSymbols := transactionRepository.GetUniqueSymbolsForPortfolio(uuid.MustParse(portfolioId))
	historicalDataPerSymbol := server.UnitOfWork.HistoricalDataRepository.GetBySymbols(uniqueSymbols)

	holdings := map[time.Time]map[string]*holding{}
	reads := make(chan readOp)
	writes := make(chan writeOp)
	go func(h map[time.Time]map[string]*holding) {
		var state = h
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.date][read.symbol]
			case write := <-writes:
				if state[write.date] == nil {
					state[write.date] = map[string]*holding{}
				}
				state[write.date][write.symbol] = write.val
				write.resp <- true
			}
		}
	}(holdings)
	var resp [][]float64
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		for _, transaction := range transactions {
			truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)
			if !truncatedTransactedAt.Equal(d) {
				continue
			}
			if holdings[truncatedTransactedAt] == nil {
				holdings[truncatedTransactedAt] = map[string]*holding{}
			}
			thisDaysHoldings := holdings[truncatedTransactedAt]
			thisDaysSymbolHoldings := thisDaysHoldings[transaction.Symbol]
			if thisDaysSymbolHoldings == nil {
				newHolding := &holding{}
				thisDaysHoldings[transaction.Symbol] = newHolding
				thisDaysSymbolHoldings = newHolding
			}
			thisDaysSymbolHoldings.amount += transaction.Amount
		}

		var dayPrice float64
		wg := sync.WaitGroup{}
		currentHoldings := holdings[d]
		c := make(chan float64, len(currentHoldings))

		for symbol, h := range currentHoldings {
			wg.Add(1)
			go func(symbol string, h *holding, c chan float64) {
				defer wg.Done()
				//var priceOfSymbolPriceAtGivenTime float64
				//for _, historicalData := range historicalDataPerSymbol[symbol] {
				//	if historicalData.Timestamp.Year() == d.Year() && historicalData.Timestamp.Month() == d.Month() && historicalData.Timestamp.Day() == d.Day() {
				//		priceOfSymbolPriceAtGivenTime = historicalData.Close
				//		break
				//	}
				//}
				symbolPriceAtGivenTime := historicalDataPerSymbol[symbol][d]
				var priceOfSymbolPriceAtGivenTime float64
				if symbolPriceAtGivenTime != nil {
					priceOfSymbolPriceAtGivenTime = symbolPriceAtGivenTime.AdjustedClose
				}
				if priceOfSymbolPriceAtGivenTime == 0 {
					read := readOp{
						date:   d.AddDate(0, 0, -1),
						symbol: symbol,
						resp:   make(chan *holding),
					}
					reads <- read
					s := <-read.resp
					if s != nil {
						priceOfSymbolPriceAtGivenTime = s.priceOfSymbolPriceAtGivenTime
					} else {
						fmt.Println("Could not find price")
					}
					// TODO Deal with symbol changes not having historical data
				}
				h.priceOfSymbolPriceAtGivenTime = priceOfSymbolPriceAtGivenTime
				h.total = priceOfSymbolPriceAtGivenTime * h.amount
				c <- h.total
				write := writeOp{
					date:   d.AddDate(0, 0, 1),
					symbol: symbol,
					val: &holding{
						amount:                        h.amount,
						priceOfSymbolPriceAtGivenTime: h.priceOfSymbolPriceAtGivenTime,
						total:                         h.total,
					},
					resp: make(chan bool),
				}
				writes <- write
				<-write.resp
			}(symbol, h, c)
		}
		wg.Wait()
		close(c)
		for elem := range c {
			dayPrice += elem
		}

		resp = append(resp, []float64{float64(d.UnixMilli()), math.Round(dayPrice*100) / 100})
	}
	allocations := entities.Allocations{}

	var amountSum float64
	for _, h := range holdings[end] {
		amountSum += h.total
	}
	for symbol, h := range holdings[end] {
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
