package entities

type Portfolio struct {
	Title        string
	Transactions Transactions `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EntityBase
}
