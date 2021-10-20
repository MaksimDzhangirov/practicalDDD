package repository

import "github.com/MaksimDzhangirov/PracticalDDD/domain/order/entity"

type OrderRepository interface {
	Create(order entity.Order) (*entity.Order, error)
}
