package transaction_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateTransactionSymbolsRequest struct {
	PortfolioId string `json:"portfolioId"`
	OldSymbol   string `json:"oldSymbol"`
	NewSymbol   string `json:"newSymbol"`
}

func UpdateTransactionSymbols(server *server.Webserver, ctx *gin.Context) {
	requestBody := UpdateTransactionSymbolsRequest{}
	_ = ctx.BindJSON(&requestBody)

	transactionRepository := server.UnitOfWork.TransactionRepository
	transactionRepository.UpdateSymbols(requestBody.PortfolioId, requestBody.OldSymbol, requestBody.NewSymbol)
	ctx.Status(http.StatusOK)
}
