package transaction_mapper

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/models/transaction_models"
)

func ToDbModel(model *transaction_models.Transaction) *entities.Transaction {
	return &entities.Transaction{
		EntityBase:      entities.EntityBase{ID: model.Id},
		TransactedAt:    model.TransactedAt,
		CurrencyCode:    model.CurrencyCode,
		TransactionType: model.TransactionType,
		Amount:          model.Amount,
		Price:           model.Price,
		Symbol:          model.Symbol,
	}
}
func ToViewModel(dbModel *entities.Transaction) *transaction_models.Transaction {
	return &transaction_models.Transaction{
		Id:              dbModel.ID,
		TransactedAt:    dbModel.TransactedAt,
		CurrencyCode:    dbModel.CurrencyCode,
		TransactionType: dbModel.TransactionType,
		Amount:          dbModel.Amount,
		Price:           dbModel.Price,
		Symbol:          dbModel.Symbol,
		Product:         dbModel.Product,
	}
}
