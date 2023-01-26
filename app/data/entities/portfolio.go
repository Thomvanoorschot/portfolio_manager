package entities

type Portfolio struct {
	EntityBase
	Title        string
	Transactions Transactions `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
