package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"math"
	"sync"
	"time"
)

type safeHoldingsPerDay struct {
	sync.RWMutex
	holdingsPerDay map[time.Time]*safeHoldings
}

func (sn *safeHoldingsPerDay) Get(d time.Time) *safeHoldings {
	sn.RLock()
	defer sn.RUnlock()
	return sn.holdingsPerDay[d]
}

type safeHoldings struct {
	sync.RWMutex
	holdings map[string]*holding
}

func (sn *safeHoldings) Add(symbol string, holding *holding) {
	sn.Lock()
	defer sn.Unlock()
	sn.holdings[symbol] = holding
}

func (sn *safeHoldings) Get(symbol string) *holding {
	sn.RLock()
	defer sn.RUnlock()
	return sn.holdings[symbol]
}

type holding struct {
	amount                 float64
	symbolPriceAtGivenTime float64
	total                  float64
}

func HistoricalDataHandler(server *server.Webserver, ctx *fasthttp.RequestCtx) {
	transactionRepository := server.UnitOfWork.TransactionRepository
	transactions := *transactionRepository.GetBuyAndSellTransactions(uuid.MustParse("f00d8e0c-d73c-411a-891e-b59cf44e8d19"))
	if len(transactions) == 0 {
		return
	}

	firstTransaction := transactions[0]
	start := helpers.TruncateToDay(firstTransaction.TransactedAt)
	end := helpers.TruncateToDay(time.Now())
	uniqueSymbols := transactionRepository.GetUniqueSymbols(uuid.MustParse("f00d8e0c-d73c-411a-891e-b59cf44e8d19"))
	historicalDataPerSymbol := server.UnitOfWork.HistoricalDataRepository.GetBySymbols(uniqueSymbols)

	holdings := safeHoldingsPerDay{
		holdingsPerDay: map[time.Time]*safeHoldings{},
	}
	holdings.holdingsPerDay[start] = &safeHoldings{
		RWMutex:  sync.RWMutex{},
		holdings: map[string]*holding{},
	}
	var resp [][]float64
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		for _, transaction := range transactions {
			truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)
			if !truncatedTransactedAt.Equal(d) {
				continue
			}
			thisDaysHoldings := holdings.holdingsPerDay[truncatedTransactedAt]
			thisDaysSymbolHoldings := thisDaysHoldings.holdings[transaction.Symbol]
			if thisDaysSymbolHoldings == nil {
				newHolding := &holding{}
				thisDaysHoldings.holdings[transaction.Symbol] = newHolding
				thisDaysSymbolHoldings = newHolding
			}
			thisDaysSymbolHoldings.amount += transaction.Amount
		}

		var dayPrice float64
		newSafeHoldings := &safeHoldings{
			RWMutex:  sync.RWMutex{},
			holdings: map[string]*holding{},
		}
		wg := sync.WaitGroup{}
		currentHoldings := holdings.holdingsPerDay[d].holdings
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
		holdings.holdingsPerDay[d.AddDate(0, 0, 1)] = newSafeHoldings

		resp = append(resp, []float64{float64(d.UnixMilli()), math.Round(dayPrice*100) / 100})
	}

	marshal, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.SetBody(marshal)
	if err != nil {
		fmt.Println(err)
		return
	}
}
