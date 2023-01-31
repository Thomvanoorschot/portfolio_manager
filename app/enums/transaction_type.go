package enums

type TransactionType string

const (
	TransactionTypeUnknown TransactionType = "Unknown"
	Purchase                               = "Purchase"
	Sale                                   = "Sale"
	Deposit                                = "Deposit"
	Withdrawal                             = "Withdrawal"
	Debit                                  = "Debit"
	Credit                                 = "Credit"
	Dividend                               = "Dividend"
	DividendTax                            = "DividendTax"
)
