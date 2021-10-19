package services

import (
	"github.com/MaksimDzhangirov/PracticalDDD/domain/bonus/repository"
	"github.com/MaksimDzhangirov/PracticalDDD/domain/services"
	domain "github.com/MaksimDzhangirov/PracticalDDD/domain/userAccount/entity"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
)

type CasinoService struct {
	bonusRepository repository.BonusRepository
	accountService  services.AccountService
	//
	// какие-нибудь другие поля
	//
}

func (s *CasinoService) Bet(account domain.Account, money value_objects.Money) error {
	bonuses, err := s.bonusRepository.FindAllEligibleFor(account, money)
	if err != nil {
		return err
	}
	//
	// какой-то код
	//
	for _, bonus := range bonuses {
		err = bonus.Apply(&account)
		if err != nil {
			return err
		}
	}
	//
	// какой-то код
	//
	err = s.accountService.Update(account)
	if err != nil {
		return err
	}
	return nil
}
