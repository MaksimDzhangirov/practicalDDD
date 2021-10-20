package services

import (
	"github.com/MaksimDzhangirov/PracticalDDD/domain/order/entity"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
)

type OrderService interface {
	Create(order entity.Order) (*entity.Order, error)
	ChangeAddress(order entity.Order, address value_objects.Address)
}
