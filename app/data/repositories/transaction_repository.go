package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func ProvideTransactionRepository(DB *gorm.DB) TransactionRepository {
	return TransactionRepository{DB: DB}
}

func (p *TransactionRepository) Create(transaction *entities.Transaction) {
	p.DB.Create(transaction)
}

func (p *TransactionRepository) GetDepositAndWithdrawalTransactions() *entities.Transactions {
	transactions := &entities.Transactions{}
	p.DB.Where("transaction_type IN ?", []entities.TransactionType{entities.Withdrawal, entities.Deposit}).Order("transacted_at asc").Find(transactions)
	return transactions
}

func (p *TransactionRepository) GetBuyAndSellTransactions() *entities.Transactions {
	transactions := &entities.Transactions{}
	p.DB.Where("transaction_type IN ?", []entities.TransactionType{entities.Buy, entities.Sell}).Order("transacted_at asc").Find(transactions)
	return transactions
}

func (p *TransactionRepository) GetUniqueSymbols() []string {
	var symbols []string
	p.DB.Model(&entities.Transaction{}).Where("symbol IS NOT NULL").Distinct("symbol").Find(&symbols)
	return symbols
}
