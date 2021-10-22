package model

import "time"

type Company struct {
	Name               string
	RegistrationNumber string
	RegistrationDate   time.Time
}
