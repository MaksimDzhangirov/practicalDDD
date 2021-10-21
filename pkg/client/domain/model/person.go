package model

import (
	"github.com/google/uuid"
	"time"
)

type Person struct {
	id uuid.UUID // локальный идентификатор
	//
	// какие-то поля
	//
	birthday time.Time
}