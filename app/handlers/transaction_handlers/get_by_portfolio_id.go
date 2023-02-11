package transaction_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	transactionMapper "github.com/Thomvanoorschot/portfolioManager/app/mappers/transaction_mappers"
	"github.com/Thomvanoorschot/portfolioManager/app/models/transaction_models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type GetByPortfolioId struct {
	transactionRepository *repositories.TransactionRepository
}

func NewGetByPortfolioId(transactionRepository *repositories.TransactionRepository) *GetByPortfolioId {
	return &GetByPortfolioId{transactionRepository: transactionRepository}
}

func (handler *GetByPortfolioId) Handle(ctx *gin.Context) {
	portfolioId := uuid.MustParse(ctx.Param("portfolioId"))

	transactions := handler.transactionRepository.GetByPortfolioId(portfolioId)
	var transactionsModel []*transaction_models.Model
	for _, transaction := range transactions {
		transactionsModel = append(transactionsModel, transactionMapper.ToViewModel(transaction))
	}
	ctx.JSON(http.StatusOK, transactionsModel)
}
