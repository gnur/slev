package slev

import (
	"testing"
	"time"
)

var opts = []SlevOpt{UseMemstore()}

func TestIntegration(t *testing.T) {
	//loop over all stores defined (currently only memstore)
	for _, opt := range opts {

		sl, err := Start(opt, MaxEvents(10), GCTime(8*time.Second))
		if err != nil {
			t.Errorf("Creating slev failed, error = %v", err)
		}
		eventsToCreate := 5
		var lastID string
		for i := 0; i < eventsToCreate; i++ {
			time.Sleep(2 * time.Millisecond)

			lastID, err = sl.NewEvent("slevtest", "simpletest", map[string]int{"i": i})
			if err != nil {
				t.Errorf("Adding event failed, error = %v", err)
			}
		}

		//there should be 10 events in the store now
		events, err := sl.store.GetEvents("", 10000)
		if err != nil {
			t.Errorf("Unable to get '', 1000 events, error = %v", err)
		}
		if len(events) != eventsToCreate {
			t.Errorf("Got too many events, expected 10, got = %v", len(events))
		}

		//there should be 10 events, only ask for 4
		events, err = sl.store.GetEvents("", 3)
		if err != nil {
			t.Errorf("Unable to get '', 4 events, error = %v", err)
		}
		if len(events) != 3 {
			t.Errorf("Got too many events, expected 4, got = %v", len(events))
		}

		//there should be no events after the lastID
		events1, err := sl.store.GetEvents(lastID, 100)
		if err != nil {
			t.Errorf("Unable to get %v, 100 events, error = %v", lastID, err)
		}
		if len(events1) != 0 {
			t.Errorf("Got too many events after %v, expected 0, got = %v", lastID, len(events1))
		}

	}

}
