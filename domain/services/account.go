package services

import "github.com/MaksimDzhangirov/PracticalDDD/domain/userAccount/entity"

// уровень предметной области
type AccountService interface {
	Update(account entity.Account) error
}