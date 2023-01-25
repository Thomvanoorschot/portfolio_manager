package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/google/uuid"
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

func (p *TransactionRepository) GetDepositAndWithdrawalTransactions(id uuid.UUID) entities.Transactions {
	transactions := entities.Transactions{}
	p.DB.Where("transaction_type IN ? AND portfolio_id = ?", []entities.TransactionType{entities.Withdrawal, entities.Deposit}, id).Order("transacted_at asc").Find(&transactions)
	return transactions
}

func (p *TransactionRepository) GetBuyAndSellTransactions(id uuid.UUID) entities.Transactions {
	transactions := entities.Transactions{}
	p.DB.Where("transaction_type IN ? AND portfolio_id = ?", []entities.TransactionType{entities.Buy, entities.Sell}, id).Order("transacted_at asc").Find(&transactions)
	return transactions
}
func (p *TransactionRepository) GetByPortfolioId(id uuid.UUID) entities.Transactions {
	transactions := entities.Transactions{}
	p.DB.Where("portfolio_id = ?", id).Order("transacted_at asc").Find(&transactions)
	return transactions
}

func (p *TransactionRepository) GetUniqueSymbolsForPortfolio(id uuid.UUID) []string {
	var symbols []string
	p.DB.Model(&entities.Transaction{}).Where("symbol IS NOT NULL AND portfolio_id = ?", id).Distinct("symbol").Find(&symbols)
	return symbols
}
func (p *TransactionRepository) GetUniqueSymbols() []string {
	var symbols []string
	p.DB.Model(&entities.Transaction{}).Where("symbol IS NOT NULL AND symbol != ?", "").Distinct("symbol").Find(&symbols)
	return symbols
}
func (p *TransactionRepository) AddToPortfolio(transactions entities.Transactions) {
	p.DB.Create(&transactions)
}
func (p *TransactionRepository) UpdateSymbols(portfolioId string, oldSymbol string, newSymbol string) {
	p.DB.Model(&entities.Transaction{}).Where("portfolio_id = ? AND symbol = ?", portfolioId, oldSymbol).Update("symbol", newSymbol)
}
