package src

import (
	"log"
	"sync"
)

// EventType , event türlerini tanımlar.
type EventType string

// Event , bir event'in temel yapısını tanımlar.
type Event struct {
	Type EventType
	Data any
}

// EventHandler , event işleyicileri için fonksiyon imzasını tanımlar.
type EventHandler func(Event)

// EventBus , event'leri yöneten ve event dinleyicilerini (subscribers) tutan yapıdır.
type EventBus struct {
	listeners map[EventType][]EventHandler
	lock      sync.Mutex
}

// NewEventBus , yeni bir EventBus örneği oluşturur.
func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[EventType][]EventHandler),
	}
}

// Publish , bir event'i yayınlar ve ilgili tüm işleyicileri tetikler.
func (bus *EventBus) Publish(event Event) {
	bus.lock.Lock()
	handlers := bus.listeners[event.Type]
	bus.lock.Unlock()

	//wg := sync.WaitGroup{}

	for _, handler := range handlers {
		log.Println("[EventBus] an event started to handle", event)
		handler(event)
		// Her bir handler'ı kendi goroutine'inde çalıştırarak asenkron işlem sağlanabilir.
		//wg.Add(1)
		//go func(wg *sync.WaitGroup, handler EventHandler) {
		//	handler(event)
		//	wg.Done()
		//	log.Println("[EventBus] an event has been handled", event)
		//}(&wg, handler)
	}

	//wg.Wait()
}

// Subscribe , belirli bir event türüne bir işleyici (handler) ekler.
func (bus *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	log.Println("[EventBus] an event listener has been registered", eventType)

	bus.listeners[eventType] = append(bus.listeners[eventType], handler)
}
