package graph_data_handlers

import (
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/Thomvanoorschot/portfolioManager/app/time_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
)

type Flag struct {
	X         int64  `json:"x"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	FillColor string `json:"fillColor"`
	LineColor string `json:"lineColor"`
}

type Trades struct {
	historicalDataRepository *repositories.HistoricalDataRepository
	transactionRepository    *repositories.TransactionRepository
}

func NewTrades(historicalDataRepository *repositories.HistoricalDataRepository,
	transactionRepository *repositories.TransactionRepository) *Trades {
	return &Trades{historicalDataRepository: historicalDataRepository,
		transactionRepository: transactionRepository}
}

func (handler *Trades) Handle(ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")
	transactions := handler.transactionRepository.GetHoldingsTransactions(uuid.MustParse(portfolioId))
	if len(transactions) == 0 {
		return
	}

	uniqueSymbols := handler.transactionRepository.GetUniqueSymbolsForPortfolio(uuid.MustParse(portfolioId))
	historicalDataPerSymbol := handler.historicalDataRepository.GetLastBySymbol(uniqueSymbols)

	var flags []*Flag
	for _, transaction := range transactions {
		truncatedTransactedAt := time_utils.TruncateToDay(transaction.TransactedAt)

		transactionTypeString := "Bought"
		transactionTitle := "B"
		filColor := "#1DA363"
		lineColor := "#15D67A"
		if transaction.TransactionType == enums.Sale {
			transactionTypeString = "Sold"
			transactionTitle = "S"
			filColor = "#AD3434"
			lineColor = "#EC1E1E"
		}

		historicalData := historicalDataPerSymbol[transaction.Symbol]
		if historicalData == nil {
			continue
		}
		// TODO Better name
		netGainOrLoss := historicalDataPerSymbol[transaction.Symbol].AdjustedClose.Sub(transaction.Price)
		gainOrLoss := decimal.NewFromInt(100).Mul(netGainOrLoss).Div(transaction.Price)
		amountPaid := transaction.Price.Mul(transaction.Amount)

		flags = append(flags, &Flag{
			X:     truncatedTransactedAt.UnixMilli(),
			Title: transactionTitle,
			Text: fmt.Sprintf("%s %.2f <b>%s</b> at %.2f for %.2f<br>Total return: %.2f%%",
				transactionTypeString,
				transaction.Amount.InexactFloat64(),
				transaction.Symbol,
				transaction.Price.InexactFloat64(),
				amountPaid.InexactFloat64(),
				gainOrLoss.InexactFloat64()),
			FillColor: filColor,
			LineColor: lineColor,
		})
	}

	ctx.JSON(http.StatusOK, flags)
}
