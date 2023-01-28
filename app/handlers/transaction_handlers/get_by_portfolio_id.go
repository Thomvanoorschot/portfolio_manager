package transaction_handlers

import (
	transaction_mapper "github.com/Thomvanoorschot/portfolioManager/app/mappers/transaction_mappers"
	"github.com/Thomvanoorschot/portfolioManager/app/models/transaction_models"
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetByPortfolioId(server *server.Webserver, ctx *gin.Context) {
	portfolioId := uuid.MustParse(ctx.Param("portfolioId"))

	transactionRepository := server.UnitOfWork.TransactionRepository
	transactions := transactionRepository.GetByPortfolioId(portfolioId)
	var transactionsModel []*transaction_models.Model
	for _, transaction := range transactions {
		transactionsModel = append(transactionsModel, transaction_mapper.ToViewModel(transaction))
	}
	ctx.JSON(http.StatusOK, transactionsModel)
}
