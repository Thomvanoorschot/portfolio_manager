package entities

type Portfolio struct {
	EntityBase
	Title        string
	Transactions Transactions `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CashBalances CashBalances `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
