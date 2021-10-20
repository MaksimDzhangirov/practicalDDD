package events

import "github.com/google/uuid"

// Интерфейс Event для описания События предметной области
type Event interface {
	Name() string
}

// Событие GeneralError
type GeneralError string

func NewGeneralError(err error) Event {
	return GeneralError(err.Error())
}

func (e GeneralError) Name() string {
	return "event.general.error"
}

// Интерфейс OrderEvent для описания События предметной области, связанных с Заказом
type OrderEvent interface {
	Event
	OrderID() uuid.UUID
}

// Событие OrderDispatched
type OrderDispatched struct {
	orderID uuid.UUID
}

func (e OrderDispatched) Name() string {
	return "event.order.dispatched"
}

func (e OrderDispatched) OrderID() uuid.UUID {
	return e.orderID
}

// Событие OrderDelivered
type OrderDelivered struct {
	orderID uuid.UUID
}

func (e OrderDelivered) Name() string {
	return "event.order.delivery.success"
}

func (e OrderDelivered) OrderID() uuid.UUID {
	return e.orderID
}

// Событие OrderDeliveryFailed
type OrderDeliveryFailed struct {
	orderID uuid.UUID
}

func (e OrderDeliveryFailed) Name() string {
	return "event.order.delivery.failed"
}

func (e OrderDeliveryFailed) OrderID() uuid.UUID {
	return e.orderID
}

type DeliveryAddressChangeFailed struct {
	orderID uuid.UUID
}

func NewDeliveryAddressChangeFailed(orderID uuid.UUID) DeliveryAddressChangeFailed {
	return DeliveryAddressChangeFailed{
		orderID: orderID,
	}
}

func (e DeliveryAddressChangeFailed) Name() string {
	return "event.order.delivery-address-change.failed"
}

func (e DeliveryAddressChangeFailed) OrderID() uuid.UUID {
	return e.orderID
}

type DeliveryAddressChanged struct {
	orderID uuid.UUID
}

func NewDeliveryAddressChanged(orderID uuid.UUID) DeliveryAddressChanged {
	return DeliveryAddressChanged{
		orderID: orderID,
	}
}

func (e DeliveryAddressChanged) Name() string {
	return "event.order.delivery-address-change.failed"
}

func (e DeliveryAddressChanged) OrderID() uuid.UUID {
	return e.orderID
}

type OrderCreated struct {
	orderID uuid.UUID
}

func NewOrderCreated(orderID uuid.UUID) OrderCreated {
	return OrderCreated{
		orderID: orderID,
	}
}

func (e OrderCreated) Name() string {
	return "event.order.dispatched"
}

func (e OrderCreated) OrderID() uuid.UUID {
	return e.orderID
}