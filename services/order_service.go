package services

import (
	"github.com/MaksimDzhangirov/PracticalDDD/domain/order/entity"
	"github.com/MaksimDzhangirov/PracticalDDD/domain/order/repository"
	"github.com/MaksimDzhangirov/PracticalDDD/events"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
)

type OrderService struct {
	repository repository.OrderRepository
	publisher events.EventPublisher
}

func (s *OrderService) Create(order entity.Order) (*entity.Order, error) {
	result, err := s.repository.Create(order)
	if err != nil {
		return nil, err
	}
	//
	// обновляем адрес в базе данных
	//
	s.publisher.Notify(events.NewOrderCreated(result.ID()))

	return result, nil
}

func (s *OrderService) ChangeAddress(order entity.Order, address value_objects.Address) {
	evt := order.ChangeAddress(address)

	s.publisher.Notify(evt) // публикуем события только внутри объект, не хранящих состояние
}