package repository

import (
	"context"
	domain "github.com/MaksimDzhangirov/PracticalDDD/domain/bankAccount/entity"
	dtopackage "github.com/MaksimDzhangirov/PracticalDDD/infrastructure/bankAccount/dto"
)

// фактическая реализация репозитория внутри инфраструктурного уровня
type BankAccountRepository struct {
	//
	// какие-то поля
	//
}

func (r *BankAccountRepository) Get(ctx context.Context, ID uint) (*domain.BankAccount, error) {
	var dto dtopackage.BankAccountGorm
	//
	// какой-то код
	//
	return &domain.BankAccount{
		ID:       dto.ID,
		IsLocked: dto.IsLocked,
		Wallet: domain.Wallet{
			Amount:   dto.Amount,
			Currency: dto.Currency.ToEntity(),
		},
		Person: dto.Person.ToEntity(),
	}, nil
}
