package dto

import (
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
	"time"
)

type PersonGorm struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	SSN       string    `gorm:"uniqueIndex;column:ssn"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Birthday  time.Time `gorm:"column:birthday"`
}

func (p *PersonGorm) ToEntity() *model.Person {
	if p == nil {
		return nil
	}

	return &model.Person{
		SSN:       p.SSN,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Birthday:  model.Birthday(p.Birthday),
	}
}