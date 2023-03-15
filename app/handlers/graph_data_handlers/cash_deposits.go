package graph_data_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/time_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type CashDeposits struct {
	transactionRepository *repositories.TransactionRepository
}

func NewCashDeposits(transactionRepository *repositories.TransactionRepository) *CashDeposits {
	return &CashDeposits{transactionRepository: transactionRepository}
}

func (handler *CashDeposits) Handle(ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	transactions := handler.transactionRepository.GetDepositAndWithdrawalTransactions(uuid.MustParse(portfolioId))
	if len(transactions) == 0 {
		return
	}

	var resp [][]float64
	firstTransaction := transactions[0]
	start := time_utils.TruncateToDay(firstTransaction.TransactedAt)
	end := time_utils.TruncateToDay(time.Now())
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		var cumulativePricePerDay decimal.Decimal
		for _, transaction := range transactions {
			truncatedTransactedAt := time_utils.TruncateToDay(transaction.TransactedAt)
			if d.After(truncatedTransactedAt) || d.Equal(truncatedTransactedAt) {
				cumulativePricePerDay = cumulativePricePerDay.Add(transaction.Price)
			}
		}
		resp = append(resp, []float64{float64(d.UnixMilli()), cumulativePricePerDay.InexactFloat64()})
	}
	ctx.JSON(http.StatusOK, resp)
}
