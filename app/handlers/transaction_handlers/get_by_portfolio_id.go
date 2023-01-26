package transaction_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type TransactionModel struct {
	Id                uuid.UUID             `json:"id"`
	TransactedAt      time.Time             `json:"transactedAt"`
	Currency          string                `json:"currency"`
	TransactionType   enums.TransactionType `json:"transactionType"`
	Product           string                `json:"product"`
	Amount            float64               `json:"amount"`
	PriceInCents      int64                 `json:"priceInCents"`
	CommissionInCents int64                 `json:"commissionInCents"`
	Symbol            string                `json:"symbol"`
}

func GetByPortfolioId(server *server.Webserver, ctx *gin.Context) {
	portfolioId := uuid.MustParse(ctx.Param("portfolioId"))

	transactionRepository := server.UnitOfWork.TransactionRepository
	transactions := transactionRepository.GetByPortfolioId(portfolioId)
	var transactionsModel []TransactionModel
	for _, transaction := range transactions {
		transactionsModel = append(transactionsModel, TransactionModel{
			Id:                transaction.ID,
			TransactedAt:      transaction.TransactedAt,
			Currency:          transaction.CurrencyCode,
			TransactionType:   transaction.TransactionType,
			Product:           transaction.Product,
			Amount:            transaction.Amount,
			PriceInCents:      transaction.PriceInCents,
			CommissionInCents: transaction.CommissionInCents,
			Symbol:            transaction.Symbol,
		})
	}
	ctx.JSON(http.StatusOK, transactionsModel)
}
