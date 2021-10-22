package model

import (
	"github.com/google/uuid"
)

type Customer struct {
	ID      uuid.UUID
	Person  *Person
	Company *Company
	Address Address
}
