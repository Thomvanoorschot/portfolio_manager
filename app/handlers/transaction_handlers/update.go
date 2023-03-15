package transaction_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/mappers/transaction_mappers"
	"github.com/Thomvanoorschot/portfolioManager/app/models/transaction_models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Update struct {
	transactionRepository *repositories.TransactionRepository
}

func NewUpdate(transactionRepository *repositories.TransactionRepository) *Update {
	return &Update{transactionRepository: transactionRepository}
}

func (handler *Update) Handle(ctx *gin.Context) {
	requestBody := &transaction_models.Transaction{}
	_ = ctx.BindJSON(requestBody)

	transaction := transaction_mapper.ToDbModel(requestBody)
	handler.transactionRepository.Update(transaction)
	ctx.Status(http.StatusOK)
}
