package transaction_mapper

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/models/transaction_models"
)

func ToDbModel(model *transaction_models.Model) *entities.Transaction {
	return &entities.Transaction{
		EntityBase:        entities.EntityBase{ID: model.Id},
		TransactedAt:      model.TransactedAt,
		CurrencyCode:      model.CurrencyCode,
		TransactionType:   model.TransactionType,
		Amount:            model.Amount,
		PriceInCents:      model.PriceInCents,
		CommissionInCents: model.CommissionInCents,
		Symbol:            model.Symbol,
	}
}
func ToViewModel(dbModel *entities.Transaction) *transaction_models.Model {
	return &transaction_models.Model{
		Id:                dbModel.ID,
		TransactedAt:      dbModel.TransactedAt,
		CurrencyCode:      dbModel.CurrencyCode,
		TransactionType:   dbModel.TransactionType,
		Amount:            dbModel.Amount,
		PriceInCents:      dbModel.PriceInCents,
		CommissionInCents: dbModel.CommissionInCents,
		Symbol:            dbModel.Symbol,
		Product:           dbModel.Product,
	}
}
