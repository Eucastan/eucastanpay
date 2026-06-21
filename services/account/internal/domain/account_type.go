package domain

type ACCType string
type AccStatus string

const (
	CurrentAccount ACCType = "current"
	SavingsAccount ACCType = "savings"
	FixedDeposit   ACCType = "fixed_deposit"
)

const (
	FreezeAccount AccStatus = "freeze"
	ActiveAccount AccStatus = "active"
	CloseAccount  AccStatus = "closed"
)
