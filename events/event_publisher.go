package events

// Интерфейс EventHandler, описывающий любой объект, который должен быть
// уведомлен о каком-либо Event
type EventHandler interface {
	Notify(event Event)
}

// EventPublisher - основная структура, уведомляющая все EventHandler
type EventPublisher struct {
	handlers map[string][]EventHandler
}

// Метод Subscribe подписывает EventHandler на определённое событие (Event)
func (e *EventPublisher) Subscribe(handler EventHandler, events ...Event) {
	for _, event := range events {
		handlers := e.handlers[event.Name()]
		handlers = append(handlers, handler)
		e.handlers[event.Name()] = handlers
	}
}

// Метод Notify уведомляет подписанный EventHandler о том, что произошло определенное событие (Event)
func (e *EventPublisher) Notify(event Event) {
	for _, handler := range e.handlers[event.Name()] {
		handler.Notify(event)
	}
}