package dto

import (
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
	"time"
)

type CompanyGorm struct {
	ID                 uint      `gorm:"primaryKey;column:id"`
	Name               string    `gorm:"column:name"`
	RegistrationNumber string    `gorm:"column:registration_number"`
	RegistrationDate   time.Time `gorm:"column:registration_date"`
}

func (c *CompanyGorm) ToEntity() *model.Company {
	if c == nil {
		return nil
	}

	return &model.Company{
		Name:               c.Name,
		RegistrationNumber: c.RegistrationNumber,
		RegistrationDate:   c.RegistrationDate,
	}
}
