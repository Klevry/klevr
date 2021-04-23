package manager

import (
	"testing"
)

func TestSendBulkEventWebHook(t *testing.T) {
	type args struct {
		url    string
		events *[]KlevrEvent
	}

	e := KlevrEvent{
		EventType: EventType(TaskCallback),
	}

	events := make([]KlevrEvent, 0)
	events = append(events, e)

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"", args{"http://127.0.0.1", &events}},
		{"nil_event", args{"http://127.0.0.1", nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendBulkEventWebHook(tt.args.url, tt.args.events)
		})
	}
}
