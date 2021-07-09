package goob

import (
	"sync"
)

// Event interface
type Event interface{}

// Pipe the Event via Write to Events. Events uses an internal buffer so it won't block Write.
// Call Stop to abort.
type Pipe struct {
	Write  func(Event)
	Events <-chan Event
	Stop   func()
}

// NewPipe instance
func NewPipe() *Pipe {
	events := make(chan Event)
	lock := sync.Mutex{}
	buf := []Event{}
	wait := make(chan struct{}, 1)
	stop := make(chan struct{})

	write := func(e Event) {
		lock.Lock()
		buf = append(buf, e)
		lock.Unlock()

		if len(wait) == 0 {
			select {
			case <-stop:
				return
			case wait <- struct{}{}:
			}
		}
	}

	go func() {
		defer close(events)

		for {
			lock.Lock()
			section := buf
			buf = []Event{}
			lock.Unlock()

			for _, e := range section {
				select {
				case <-stop:
					return
				case events <- e:
				}
			}

			select {
			case <-stop:
				return
			case <-wait:
			}
		}
	}()

	return &Pipe{write, events, func() { close(stop) }}
}
