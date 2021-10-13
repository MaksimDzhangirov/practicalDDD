package domain

import (
	"errors"
	"time"
)

// Сущность внутри уровня предметной области
type BankAccount struct {
	ID       int
	IsLocked bool
	Wallet   Wallet
	Person   Person
}

func (ba *BankAccount) Add(other Wallet) error {
	if ba.IsLocked {
		return errors.New("account is locked")
	}
	//
	// что-то делаем
	//
	return nil
}

// правильно - сущность BankAccount проверяет свои собственные инварианты
func (ba *BankAccount) Deduct(other Wallet) error {
	if ba.IsLocked {
		return errors.New("account is locked")
	}

	result, err := ba.Wallet.Deduct(other)
	if err != nil {
		return err
	}

	ba.Wallet = result

	return nil
}

type Currency struct {
	ID       uint
	Code     string
	Name     string
	HtmlCode string
}

type Person struct {
	ID          uint
	FirstName   string
	LastName    string
	DateOfBirth time.Time
}

type Wallet struct {
	Amount   int
	Currency Currency
}

func (c Currency) IsEqual(other Currency) bool {
	return other.ID == c.ID
}

// правильно - объект-значение Wallet проверяет свои собственные инварианты
func (w Wallet) Deduct(other Wallet) (Wallet, error) {
	if !other.Currency.IsEqual(w.Currency) {
		return Wallet{}, errors.New("currencies must be the same")
	}
	if other.Amount > w.Amount {
		return Wallet{}, errors.New("insufficient funds")
	}

	return Wallet{
		Amount:   w.Amount - other.Amount,
		Currency: w.Currency,
	}, nil
}
