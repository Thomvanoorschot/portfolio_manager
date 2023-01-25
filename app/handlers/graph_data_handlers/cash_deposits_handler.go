package graph_data_handlers

import (
	"github.com/Rhymond/go-money"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func CashDepositsHandler(server *server.Webserver, ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	transactions := server.UnitOfWork.TransactionRepository.GetDepositAndWithdrawalTransactions(uuid.MustParse(portfolioId))
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
	ctx.JSON(http.StatusOK, resp)
}