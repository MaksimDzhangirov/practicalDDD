package model

import "github.com/google/uuid"

// Сущность
type BankAccount struct {
	id       uuid.UUID
	iban     string
	amount   int
	currency Currency
}

func NewBankAccount(currency Currency) BankAccount {
	return BankAccount{
		//
		// определяем поля
		//
	}
}

func (ba BankAccount) HasMoney() bool {
	return ba.amount > 0
}

func (ba BankAccount) InDebt() bool {
	return ba.amount < 0
}

func (ba BankAccount) IsForCurrency(currency Currency) bool {
	return ba.currency.Equal(currency)
}

type BankAccounts []BankAccount

func (bas BankAccounts) HasMoney() bool {
	for _, ba := range bas {
		if ba.HasMoney() {
			return true
		}
	}

	return false
}

func (bas BankAccounts) InDebt() bool {
	for _, ba := range bas {
		if ba.InDebt() {
			return true
		}
	}

	return false
}

func (bas BankAccounts) HasCurrency(currency Currency) bool {
	for _, ba := range bas {
		if ba.IsForCurrency(currency) {
			return true
		}
	}

	return false
}

func (bas BankAccounts) AddMoney(amount int, currency Currency) error {
	panic("not implemented")
	return nil
}
