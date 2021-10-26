package slev

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Event struct {
	ID     string    `storm:"unique"`
	Source string    `storm:"index"`
	Type   string    `storm:"index"`
	Time   time.Time `storm:"index"`
	Data   interface{}
}

type Store interface {
	//AddEvent adds an event to the store
	AddEvent(Event) error
	//GetEvents returns all the events after the optional "after" id, with a max of count
	GetEvents(after string, count int) ([]Event, error)
	//GetEvents returns all the events after the optional "after" id, with a max of count
	GetEvent(id string) (Event, error)
	//Cleanup removes all events over the max, keeping only the newest events
	Cleanup(max int) error
}

type Slev struct {
	max         int
	gcInterval  time.Duration
	lastCleanup time.Time
	store       Store
	host        *http.Server
}

type SlevOpt func(*Slev) error

func UseDefaultHTTPServer(listenAddress string) SlevOpt {
	return func(s *Slev) error {
		go func() {
			http.HandleFunc("/events", func(w http.ResponseWriter, req *http.Request) {
				after := req.URL.Query().Get("after")
				limit := req.URL.Query().Get("limit")

				if after != "" {
					return
				}
				if limit != "" {
					return
				}
				events := s.RawEvents()
				js, err := json.Marshal(events)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			})
			http.ListenAndServe(listenAddress, nil)
		}()
		return nil
	}
}

func MaxEvents(n int) SlevOpt {
	return func(s *Slev) error {
		s.max = n
		return nil
	}
}

func UseMemstore() SlevOpt {
	return func(s *Slev) error {
		s.store = NewMemstore()
		return nil
	}
}

func GCTime(n time.Duration) SlevOpt {
	return func(s *Slev) error {
		if n < time.Second {
			return errors.New("Invalid duration, should be longer then 1 second")
		}
		s.gcInterval = n
		return nil
	}
}

func Start(opts ...SlevOpt) (*Slev, error) {
	s := Slev{
		max:        500,
		gcInterval: 63 * time.Second,
	}
	var err error

	for _, opt := range opts {
		err = opt(&s)
		if err != nil {
			return nil, err
		}
	}
	if s.store == nil {
		s.store = NewMemstore()
	}

	go s.cleanupLoop()

	return &s, nil
}

func (s *Slev) cleanupLoop() {
	for {
		time.Sleep(s.gcInterval)
		s.store.Cleanup(s.max)
	}
}

func (s *Slev) NewEvent(source, typ string, data interface{}) (string, error) {

	id, err := NewID()
	if err != nil {
		return "", err
	}
	e := Event{
		ID:     id,
		Source: source,
		Type:   typ,
		Time:   time.Now(),
		Data:   data,
	}
	err = s.store.AddEvent(e)

	return e.ID, err
}

func (s *Slev) RawEvents() []Event {
	e, _ := s.store.GetEvents("", 10000)

	return e
}
