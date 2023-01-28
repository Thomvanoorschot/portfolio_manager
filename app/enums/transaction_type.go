package enums

type TransactionType string

const (
	TransactionTypeUnknown TransactionType = "Unknown"
	Purchase                               = "Purchase"
	Sale                                   = "Sale"
	Deposit                                = "Deposit"
	Withdrawal                             = "Withdrawal"
)
