package transaction_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Update(server *server.Webserver, ctx *gin.Context) {
	requestBody := TransactionModel{}
	_ = ctx.BindJSON(&requestBody)

	transactionRepository := server.UnitOfWork.TransactionRepository
	transactionRepository.Update(&entities.Transaction{
		EntityBase:        entities.EntityBase{ID: requestBody.Id},
		TransactedAt:      requestBody.TransactedAt,
		CurrencyCode:      requestBody.Currency,
		TransactionType:   requestBody.TransactionType,
		Amount:            requestBody.Amount,
		PriceInCents:      requestBody.PriceInCents,
		CommissionInCents: requestBody.CommissionInCents,
		Symbol:            requestBody.Symbol,
	})
	ctx.Status(http.StatusOK)
}
