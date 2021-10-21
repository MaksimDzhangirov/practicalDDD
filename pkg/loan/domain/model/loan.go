package model

import (
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
	"github.com/google/uuid"
)

const (
	LongTerm = iota
	ShortTerm
)

type Loan struct {
	ID                    uuid.UUID
	Type                  int
	BankAccountID         uuid.UUID
	Amount                value_objects.Money
	RequiredLifeInsurance bool
}

type LoanFactory struct{}

func (f *LoanFactory) CreateShortTermLoan(bankAccountID uuid.UUID, amount value_objects.Money) Loan {
	return Loan{
		Type:          ShortTerm,
		BankAccountID: bankAccountID,
		Amount:        amount,
	}
}

func (f *LoanFactory) CreateLongTermLoan(bankAccountID uuid.UUID, amount value_objects.Money) Loan {
	return Loan{
		Type:                  LongTerm,
		BankAccountID:         bankAccountID,
		Amount:                amount,
		RequiredLifeInsurance: true,
	}
}
