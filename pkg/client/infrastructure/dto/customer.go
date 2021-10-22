package dto

import (
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
	"github.com/google/uuid"
)

type CustomerGorm struct {
	ID        uint         `gorm:"primaryKey;column:id"`
	UUID      string       `gorm:"uniqueIndex;column:id"`
	PersonID  uint         `gorm:"column:person_id"`
	Person    *PersonGorm  `gorm:"foreignKey:PersonID"`
	CompanyID uint         `gorm:"column:company_id"`
	Company   *CompanyGorm `gorm:"foreignKey:CompanyID"`
	Street    string       `gorm:"column:street"`
	Number    string       `gorm:"column:number"`
	Postcode  string       `gorm:"column:postcode"`
	City      string       `gorm:"column:city"`
}

func (c CustomerGorm) ToEntity() (model.Customer, error) {
	parsed, err := uuid.Parse(c.UUID)
	if err != nil {
		return model.Customer{}, err
	}

	return model.Customer{
		ID:      parsed,
		Person:  c.Person.ToEntity(),
		Company: c.Company.ToEntity(),
		Address: model.Address{
			Street:   c.Street,
			Number:   c.Number,
			Postcode: c.Postcode,
			City:     c.City,
		},
	}, nil
}

type CustomerJSON struct {

}

func (c CustomerJSON) ToEntity() (model.Customer, error) {
	panic("not implemented")
	return model.Customer{}, nil
}