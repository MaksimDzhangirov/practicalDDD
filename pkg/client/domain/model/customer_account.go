package model

import "github.com/google/uuid"

type CustomerAccount struct {
	id uuid.UUID // глобальный идентификатор
	person  *Person
	company *Company
	//
	// какие-то поля
	//
}