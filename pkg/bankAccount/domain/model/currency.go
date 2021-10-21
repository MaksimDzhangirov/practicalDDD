package model

import "github.com/google/uuid"

// Сущность
type Currency struct {
	id uuid.UUID
	//
	// какие-то поля
	//
}

func (c Currency) Equal(other Currency) bool {
	return c.id == other.id
}
