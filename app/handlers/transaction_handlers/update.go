package transaction_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/mappers/transaction_mappers"
	"github.com/Thomvanoorschot/portfolioManager/app/models/transaction_models"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Update(server *server.Webserver, ctx *gin.Context) {
	requestBody := &transaction_models.Model{}
	_ = ctx.BindJSON(requestBody)

	transaction := transaction_mapper.ToDbModel(requestBody)
	transactionRepository := server.UnitOfWork.TransactionRepository
	transactionRepository.Update(transaction)
	ctx.Status(http.StatusOK)
}
