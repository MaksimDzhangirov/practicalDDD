package model

import (
	"errors"
	"github.com/google/uuid"
)

// Сущность и агрегат
type CustomerAccount struct {
	id        uuid.UUID
	isDeleted bool
	isLocked  bool
	//
	// какие-то поля
	//
	accounts BankAccounts
	//
	// какие-то поля
	//
}

func (ca *CustomerAccount) GetIBANForCurrency(currency Currency) (string, error) {
	for _, account := range ca.accounts {
		if account.IsForCurrency(currency) {
			return account.iban, nil
		}
	}
	return "", errors.New("this account does not support this currency")
}

func (ca *CustomerAccount) MarkAsDeleted() error {
	if ca.accounts.HasMoney() {
		return errors.New("there are still money on bank account")
	}
	if ca.accounts.InDebt() {
		return errors.New("bank account is in debt")
	}

	ca.isDeleted = true

	return nil
}

func (ca *CustomerAccount) CreateAccountForCurrency(currency Currency) error {
	if ca.accounts.HasCurrency(currency) {
		return errors.New("there is already bank account for that currency")
	}
	ca.accounts = append(ca.accounts, NewBankAccount(currency))

	return nil
}

func (ca *CustomerAccount) AddMoney(amount int, currency Currency) error {
	if ca.isDeleted {
		return errors.New("account is deleted")
	}
	if ca.isLocked {
		return errors.New("account is locked")
	}

	return ca.accounts.AddMoney(amount, currency)
}
