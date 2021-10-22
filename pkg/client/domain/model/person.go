package model

import (
	"time"
)

type Birthday time.Time

type Person struct {
	SSN       string
	FirstName string
	LastName  string
	Birthday  Birthday
}
