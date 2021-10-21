package model

import "github.com/google/uuid"

type Company struct {
	id uuid.UUID // локальный идентификатор
	//
	// какие-то поля
	//
	isLiquid bool
}
