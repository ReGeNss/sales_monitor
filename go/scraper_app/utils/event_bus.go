package utils

import "reflect"
	
type EventBus interface {
	Publish(event interface{})
	Subscribe(event interface{}, handler func(payload interface{})) error
}

type eventBusImpl struct {
	subscribers map[reflect.Type][]func(payload interface{})
}

func NewEventBus() EventBus {
	return &eventBusImpl{
		subscribers: make(map[reflect.Type][]func(payload interface{})),
	}
}

func (b *eventBusImpl) Publish(event interface{}) {
	eventName := reflect.TypeOf(event)

	if handlers, exist := b.subscribers[eventName]; exist {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

func (b *eventBusImpl) Subscribe(event interface{}, handler func(payload interface{})) error {
	eventName := reflect.TypeOf(event)
	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
	return nil
}
