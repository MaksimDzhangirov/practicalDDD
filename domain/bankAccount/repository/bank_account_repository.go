package repository

import (
	"context"
	domain "github.com/MaksimDzhangirov/PracticalDDD/domain/bankAccount/entity"
)

// Интерфейс репозитория внутри уровня предметной области
type BankAccountRepository interface {
	Get(ctx context.Context, ID int) (*domain.BankAccount, error)
}
