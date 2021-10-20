package events

import (
	"fmt"
	"github.com/google/uuid"
)

type EmailEvent interface {
	Event
	EmailID() uuid.UUID
}

type EmailSent struct {
	emailID uuid.UUID
}

func (e EmailSent) Name() string {
	return "event.email.sent"
}

func (e EmailSent) EmailID() uuid.UUID {
	return e.emailID
}

type EmailHandler struct {
	//
	// какие-то поля
	//
}

func (e *EmailHandler) Notify(event Event) {
	switch actualEvent := event.(type) {
	case EmailSent:
		fmt.Println(actualEvent)
		//
		// что-то делаем
		//
	default:
		return
	}
}
