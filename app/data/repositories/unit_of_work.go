package repositories

type UnitOfWork struct {
	TransactionRepository    *TransactionRepository
	PortfolioRepository      *PortfolioRepository
	HistoricalDataRepository *HistoricalDataRepository
	AllocationRepository     *AllocationRepository
	CashBalanceRepository    *CashBalanceRepository
}
