package model

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	id      uuid.UUID
	person  *Person
	company *Company
	//
	// какие-то поля
	//
}

func (c *Customer) IsLegal() bool {
	if c.person != nil {
		return c.person.birthday.AddDate(18, 0, 0).Before(time.Now())
	} else {
		return c.company.isLiquid
	}
}
