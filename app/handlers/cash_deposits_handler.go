package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"time"
)

func CashDepositsHandler(server *server.Webserver, ctx *fasthttp.RequestCtx) {
	transactions := *server.UnitOfWork.TransactionRepository.GetDepositAndWithdrawalTransactions(uuid.MustParse("58ade236-bdec-444b-a98e-653f8a4eabc3"))
	if len(transactions) == 0 {
		return
	}

	var resp [][]float64
	firstTransaction := transactions[0]
	start := helpers.TruncateToDay(firstTransaction.TransactedAt)
	end := helpers.TruncateToDay(time.Now())
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		var cumulativePriceInCentsPerDay int64
		for _, transaction := range transactions {
			truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)
			if d.After(truncatedTransactedAt) || d.Equal(truncatedTransactedAt) {
				cumulativePriceInCentsPerDay += transaction.PriceInCents
			}
		}
		resp = append(resp, []float64{float64(d.UnixMilli()), money.New(cumulativePriceInCentsPerDay, firstTransaction.CurrencyCode).AsMajorUnits()})
	}
	marshal, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.SetBody(marshal)
	ctx.SetContentType("application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		fmt.Println(err)
		return
	}
}
