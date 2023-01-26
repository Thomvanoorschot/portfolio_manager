package enums

type TransactionType int

const (
	Unknown TransactionType = iota + 1
	Buy
	Sell
	Deposit
	Withdrawal
)
