package repositories

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/entities"
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: DB}
}

func (p *TransactionRepository) Create(transaction *entities.Transaction) {
	p.DB.Create(transaction)
}

func (p *TransactionRepository) GetDepositAndWithdrawalTransactions(id uuid.UUID) entities.Transactions {
	transactions := entities.Transactions{}
	p.DB.Where("transaction_type IN ? AND portfolio_id = ?", []enums.TransactionType{enums.Withdrawal, enums.Deposit}, id).Order("transacted_at asc").Find(&transactions)
	return transactions
}

func (p *TransactionRepository) GetHoldingsTransactionsPerSymbol(id uuid.UUID, symbol string) entities.Transactions {
	transactions := entities.Transactions{}
	p.DB.Where("transaction_type IN ? "+
		"AND portfolio_id = ? "+
		"AND symbol != ? "+
		"AND symbol = ?",
		[]enums.TransactionType{enums.Purchase, enums.Sale},
		id,
		"",
		symbol).Order("transacted_at asc").Find(&transactions)
	return transactions
}
func (p *TransactionRepository) GetHoldingsTransactions(id uuid.UUID) entities.Transactions {
	transactions := entities.Transactions{}
	p.DB.Where("transaction_type IN ? AND portfolio_id = ? AND symbol != ?", []enums.TransactionType{enums.Purchase, enums.Sale}, id, "").Order("transacted_at asc").Find(&transactions)
	return transactions
}
func (p *TransactionRepository) GetByPortfolioId(id uuid.UUID) entities.Transactions {
	transactions := entities.Transactions{}
	p.DB.Where("portfolio_id = ?", id).Order("transacted_at asc").Find(&transactions)
	return transactions
}
func (p *TransactionRepository) GetLastTransaction() *entities.Transaction {
	transaction := &entities.Transaction{}
	p.DB.Order("transacted_at desc").First(transaction)
	return transaction
}

func (p *TransactionRepository) GetUniqueSymbolsForPortfolio(id uuid.UUID) []string {
	var symbols []string
	p.DB.Model(&entities.Transaction{}).Where("symbol IS NOT NULL AND portfolio_id = ? AND symbol != ? AND transaction_type IN ?", id, "", []enums.TransactionType{enums.Purchase, enums.Sale}).Distinct("symbol").Find(&symbols)
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

func (p *TransactionRepository) Update(transaction *entities.Transaction) {
	p.DB.Model(&entities.Transaction{}).Where("id = ?", transaction.ID).Select("transacted_at",
		"currency_code",
		"transaction_type",
		"amount",
		"price_in_cents",
		"commission_in_cents",
		"symbol").Updates(transaction)
}

//
//TransactedAt:      requestBody.TransactedAt,
//CurrencyCode:      requestBody.CurrencyCode,
//TransactionType:   requestBody.TransactionType,
//Amount:            requestBody.Amount,
//PriceInCents:      requestBody.PriceInCents,
//CommissionInCents: requestBody.CommissionInCents,
//Symbol:            requestBody.Symbol,
