package repository

import (
	"github.com/MaksimDzhangirov/PracticalDDD/domain/bonus/entity"
	domain "github.com/MaksimDzhangirov/PracticalDDD/domain/userAccount/entity"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
)

type BonusRepository struct {

}

func (r *BonusRepository) FindAllEligibleFor(account domain.Account, money value_objects.Money) ([]entity.Bonus, error) {
	var result []entity.Bonus
	//
	// какой-то код
	//
	return result, nil
}