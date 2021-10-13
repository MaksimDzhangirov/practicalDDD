package dto

import domain "github.com/MaksimDzhangirov/PracticalDDD/domain/bankAccount/entity"

// DTO внутри инфраструктурного уровня
type BankAccountGorm struct {
	ID         int          `gorm:"primaryKey";column:id`
	IsLocked   bool         `gorm:"column:is_locked"`
	Amount     int          `gorm:"column:amount"`
	CurrencyID uint         `gorm:"column:currency_id"`
	Currency   CurrencyGorm `gorm:"foreignKey:CurrencyID"`
	PersonID   uint         `gorm:"column:person_id"`
	Person     PersonGorm   `gorm:"foreignKey:PersonID"`
}

type CurrencyGorm struct {

}

type PersonGorm struct {

}

func (cg *CurrencyGorm) ToEntity() domain.Currency {
	return domain.Currency{}
}

func (pg *PersonGorm) ToEntity() domain.Person {
	return domain.Person{}
}