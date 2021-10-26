package slev

import (
	"sync"
)

type Memstore struct {
	lock   sync.Mutex
	events []Event
}

func NewMemstore() *Memstore {
	m := Memstore{}
	return &m
}

func (m *Memstore) AddEvent(e Event) error {
	m.lock.Lock()
	m.events = append(m.events, e)
	m.lock.Unlock()
	return nil
}

func (m *Memstore) GetEvents(after string, count int) ([]Event, error) {
	var events []Event
	for _, e := range m.events {
		if e.ID > after {
			events = append(events, e)
			if len(events) == count {
				break
			}
		}
	}

	return events, nil
}

func (m *Memstore) GetEvent(id string) (Event, error) {
	for _, e := range m.events {
		if e.ID == id {
			return e, nil
		}
	}
	return Event{}, nil
}

func (m *Memstore) Cleanup(max int) error {
	m.lock.Lock()

	if len(m.events) > max {
		start := len(m.events) - max
		m.events = m.events[start : start+max]
	}

	m.lock.Unlock()
	return nil
}
