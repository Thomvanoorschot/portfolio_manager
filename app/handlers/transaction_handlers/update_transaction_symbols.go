package transaction_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateTransactionSymbolsRequest struct {
	PortfolioId string         `json:"portfolioId"`
	SymbolPairs []SymbolUpdate `json:"symbolPairs"`
}

type SymbolUpdate struct {
	OldSymbol string `json:"oldSymbol"`
	NewSymbol string `json:"newSymbol"`
}

type UpdateTransactionSymbols struct {
	transactionRepository *repositories.TransactionRepository
}

func NewUpdateTransactionSymbols(transactionRepository *repositories.TransactionRepository) *UpdateTransactionSymbols {
	return &UpdateTransactionSymbols{transactionRepository: transactionRepository}
}

func (handler *UpdateTransactionSymbols) Handle(ctx *gin.Context) {
	requestBody := UpdateTransactionSymbolsRequest{}
	_ = ctx.BindJSON(&requestBody)
	for _, symbolPair := range requestBody.SymbolPairs {
		handler.transactionRepository.UpdateSymbols(requestBody.PortfolioId, symbolPair.OldSymbol, symbolPair.NewSymbol)
	}
	ctx.Status(http.StatusOK)
}
