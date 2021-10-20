package entity

import (
	"github.com/MaksimDzhangirov/PracticalDDD/events"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
	"github.com/google/uuid"
)

type Order struct {
	id uuid.UUID
	//
	// какие-то поля
	//
	isDispatched bool
	deliverAddress value_objects.Address
}

func (o Order) ID() uuid.UUID {
	return o.id
}

func (o Order) ChangeAddress(address value_objects.Address) events.Event {
	if o.isDispatched {
		return events.NewDeliveryAddressChangeFailed(o.ID())
	}
	//
	// какой-то код
	//
	return events.NewDeliveryAddressChanged(o.ID())
}