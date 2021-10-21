package model

import (
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
	"github.com/google/uuid"
)

type Investment interface {
	Amount() value_objects.Money
}

type EtfInvestment struct {
	ID             uuid.UUID
	EtfID          uuid.UUID
	InvestedAmount value_objects.Money
	BankAccountID  uuid.UUID
}

func (e EtfInvestment) Amount() value_objects.Money {
	return e.InvestedAmount
}

type StockInvestment struct {
	ID             uuid.UUID
	CompanyID      uuid.UUID
	InvestedAmount value_objects.Money
	BankAccountID  uuid.UUID
}

func (s StockInvestment) Amount() value_objects.Money {
	return s.InvestedAmount
}

type CryptoInvestment struct {
	ID               uuid.UUID
	CryptoCurrencyID uuid.UUID
	InvestedMoney    value_objects.Money
	BankAccountID    uuid.UUID
}

func (c CryptoInvestment) Amount() value_objects.Money {
	return c.InvestedMoney
}

type InvestmentSpecification interface {
	Amount() value_objects.Money
	BankAccountID() uuid.UUID
	TargetID() uuid.UUID
}

type InvestmentFactory interface {
	Create(specification InvestmentSpecification) Investment
}

type EtfInvestmentFactory struct{}

func (f *EtfInvestmentFactory) Create(specification InvestmentSpecification) Investment {
	return EtfInvestment{
		EtfID:          specification.TargetID(),
		InvestedAmount: specification.Amount(),
		BankAccountID:  specification.BankAccountID(),
	}
}

type StockInvestmentFactory struct{}

func (f *StockInvestmentFactory) Create(specification InvestmentSpecification) Investment {
	return StockInvestment{
		CompanyID:      specification.TargetID(),
		InvestedAmount: specification.Amount(),
		BankAccountID:  specification.BankAccountID(),
	}
}
