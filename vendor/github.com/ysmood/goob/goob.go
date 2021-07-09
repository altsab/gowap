package goob

import (
	"sync"
)

// Observable hub
type Observable struct {
	lock        *sync.Mutex
	subscribers map[Subscriber]*Pipe
}

// Subscriber type
type Subscriber <-chan Event

// New observable instance
func New() *Observable {
	ob := &Observable{
		lock:        &sync.Mutex{},
		subscribers: map[Subscriber]*Pipe{},
	}
	return ob
}

// Publish message to the queue
func (ob *Observable) Publish(e Event) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	for _, p := range ob.subscribers {
		p.Write(e)
	}
}

// Subscribe message
func (ob *Observable) Subscribe() Subscriber {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	p := NewPipe()

	if ob.subscribers == nil {
		p.Stop()
	} else {
		ob.subscribers[p.Events] = p
	}

	return p.Events
}

// Unsubscribe from observable
func (ob *Observable) Unsubscribe(s Subscriber) {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	if p, has := ob.subscribers[s]; has {
		p.Stop()
		delete(ob.subscribers, s)
	}
}

// Close subscribers
func (ob *Observable) Close() {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	for _, p := range ob.subscribers {
		p.Stop()
	}

	ob.subscribers = nil
}

// Len of the subscribers
func (ob *Observable) Len() int {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	return len(ob.subscribers)
}
