package services

import (
	domain "github.com/MaksimDzhangirov/PracticalDDD/domain/bankAccount/entity"
	"github.com/MaksimDzhangirov/PracticalDDD/domain/exchangeRate/repository"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
)

type ExchangeRateService interface {
	IsConversionPossible (from domain.Currency, to domain.Currency) bool
	Convert(to domain.Currency, from value_objects.Money) (value_objects.Money, error)
}

type DefaultExchangeRateService struct {
	repository repository.ExchangeRateRepository
}

func NewExchangeRateService(repository repository.ExchangeRateRepository) ExchangeRateService {
	return &DefaultExchangeRateService{
		repository: repository,
	}
}

func (s *DefaultExchangeRateService) IsConversionPossible(from domain.Currency, to domain.Currency) bool {
	var result bool
	//
	// какой-то код
	//
	return result
}

func (s *DefaultExchangeRateService) Convert(to domain.Currency, from value_objects.Money) (value_objects.Money, error) {
	var result value_objects.Money
	//
	// какой-то код
	//
	return result, nil
}