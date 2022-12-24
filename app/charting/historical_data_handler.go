package charting

import (
	"encoding/json"
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/Thomvanoorschot/portfolioManager/app/infrastructure"
	"net/http"
	"time"
)

type holding struct {
	amount                 float64
	symbolPriceAtGivenTime float64
	total                  float64
}

func HistoricalDataHandler(server *infrastructure.Server, _ *http.Request, response http.ResponseWriter) {
	transactionRepository := server.UnitOfWork.TransactionRepository
	transactions := *transactionRepository.GetBuyAndSellTransactions()
	if len(transactions) == 0 {
		return
	}

	firstTransaction := transactions[0]
	start := helpers.TruncateToDay(firstTransaction.TransactedAt)
	end := helpers.TruncateToDay(time.Now())
	uniqueSymbols := transactionRepository.GetUniqueSymbols()
	historicalDataPerSymbol := server.UnitOfWork.HistoricalDataRepository.GetBySymbols(uniqueSymbols)

	holdings := map[time.Time]map[string]*holding{}
	holdings[start] = map[string]*holding{}
	var resp [][]float64
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		for _, transaction := range transactions {
			truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)
			if !truncatedTransactedAt.Equal(d) {
				continue
			}

			if holdings[truncatedTransactedAt][transaction.Symbol] == nil {
				holdings[truncatedTransactedAt][transaction.Symbol] = &holding{}
			}
			holdings[truncatedTransactedAt][transaction.Symbol].amount += transaction.Amount
		}

		var dayPrice float64
		holdings[d.AddDate(0, 0, 1)] = map[string]*holding{}
		for symbol, h := range holdings[d] {
			holdings[d.AddDate(0, 0, 1)][symbol] = &holding{
				amount:                 h.amount,
				symbolPriceAtGivenTime: h.symbolPriceAtGivenTime,
				total:                  h.total,
			}
			var symbolPriceAtGivenTime float64
			for _, historicalData := range historicalDataPerSymbol[symbol] {
				if historicalData.Timestamp.Year() == d.Year() && historicalData.Timestamp.Month() == d.Month() && historicalData.Timestamp.Day() == d.Day() {
					symbolPriceAtGivenTime = historicalData.Close
					break
				}
			}
			if symbolPriceAtGivenTime == 0 {
				previousDaysHolding := holdings[d.AddDate(0, 0, -1)]
				s := previousDaysHolding[symbol]
				if s != nil {
					symbolPriceAtGivenTime = s.symbolPriceAtGivenTime
				}
				// TODO Deal with symbol changes not having historical data
			}
			h.symbolPriceAtGivenTime = symbolPriceAtGivenTime
			h.total = symbolPriceAtGivenTime * h.amount
			dayPrice += h.total
		}

		resp = append(resp, []float64{float64(d.UnixMilli()), dayPrice})
	}

	marshal, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = response.Write(marshal)
	if err != nil {
		fmt.Println(err)
		return
	}
}
