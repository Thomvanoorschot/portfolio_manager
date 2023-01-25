package graph_data_handlers

import (
	"fmt"
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/helpers"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Flag struct {
	X         int64  `json:"x"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	FillColor string `json:"fillColor"`
	LineColor string `json:"lineColor"`
}

func TradesHandler(server *server.Webserver, ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")
	transactionRepository := server.UnitOfWork.TransactionRepository
	transactions := transactionRepository.GetBuyAndSellTransactions(uuid.MustParse(portfolioId))
	if len(transactions) == 0 {
		return
	}

	uniqueSymbols := transactionRepository.GetUniqueSymbolsForPortfolio(uuid.MustParse(portfolioId))
	historicalDataPerSymbol := server.UnitOfWork.HistoricalDataRepository.GetLastBySymbol(uniqueSymbols)

	var flags []*Flag
	for _, transaction := range transactions {
		truncatedTransactedAt := helpers.TruncateToDay(transaction.TransactedAt)

		transactionTypeString := "Bought"
		transactionTitle := "B"
		filColor := "#1DA363"
		lineColor := "#15D67A"
		if transaction.TransactionType == entities.Sell {
			transactionTypeString = "Sold"
			transactionTitle = "S"
			filColor = "#AD3434"
			lineColor = "#EC1E1E"
		}

		historicalData := historicalDataPerSymbol[transaction.Symbol]
		if historicalData == nil {
			continue
		}
		gainOrLoss := 100 * (historicalDataPerSymbol[transaction.Symbol].AdjustedClose*100 - float64(transaction.PriceInCents)) / float64(transaction.PriceInCents)
		flags = append(flags, &Flag{
			X:         truncatedTransactedAt.UnixMilli(),
			Title:     transactionTitle,
			Text:      fmt.Sprintf("%s %.2f <b>%s</b> at %.2f for %.2f<br>Total return: %.2f%%", transactionTypeString, transaction.Amount, transaction.Symbol, float64(transaction.PriceInCents)/100, float64(transaction.PriceInCents)*transaction.Amount/100, gainOrLoss),
			FillColor: filColor,
			LineColor: lineColor,
		})
	}

	ctx.JSON(http.StatusOK, flags)
}
