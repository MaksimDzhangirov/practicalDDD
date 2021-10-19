package services

import (
	"fmt"
	"github.com/MaksimDzhangirov/PracticalDDD/domain/bonus/repository"
	"github.com/MaksimDzhangirov/PracticalDDD/domain/userAccount/entity"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
)

// правильно - состояние передаётся в виде аргумента current
type TransactionService struct {
	bonusRepository repository.BonusRepository
}

func (s *TransactionService) Deposit(current value_objects.Money, account entity.Account, money value_objects.Money) (value_objects.Money, error) {
	bonuses, err := s.bonusRepository.FindAllEligibleFor(account, money)
	if err != nil {
		return value_objects.Money{}, err
	}
	fmt.Printf("%v", bonuses)
	//
	// какой-то код
	//
	return current.Add(money) // возвращаем новое значение, которое представляет новое состояние
}
