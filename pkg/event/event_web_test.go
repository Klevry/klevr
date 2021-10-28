package event

import (
	"fmt"
	"testing"
	"time"

	"github.com/Klevry/klevr/pkg/queue"
)

func TestNewEventWeb(t *testing.T) {
	type args struct {
		opt KlevrEventOption
	}
	tests := []struct {
		name string
		args args
		want EventManager
	}{
		// TODO: Add test cases.
		{
			name: "neweventweb",
			args: args{opt: KlevrEventOption{URL: []string{"http://localhost", "http://127.0.0.1"}, Web_HookCount: 10, Web_HookTerm: 3}},
			want: &EventWeb{eventQueue: queue.NewMutexQueue(), url: []string{"http://localhost", "http://127.0.0.1"}, hookCount: 10, hookTerm: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEventWeb(tt.args.opt)
			if got == EventManager(tt.want) {
				t.Errorf("NewEventWeb() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestEventWeb_AddEvent(t *testing.T) {
	type args struct {
		event *KlevrEvent
	}
	tests := []struct {
		name string
		e    EventManager
		args args
	}{
		// TODO: Add test cases.

		{
			name: "event object call addevent()",
			e:    NewEventWeb(KlevrEventOption{URL: []string{"http://localhost:9985/event"}, Web_HookCount: 2, Web_HookTerm: 5}),
			args: args{event: &KlevrEvent{EventType: AgentConnect, AgentKey: "test-key", GroupID: uint64(123), Result: "aaa", Log: "log-agdfads"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// hookCount test
			tt.e.AddEvent(tt.args.event)
			event1 := &KlevrEvent{EventType: AgentConnect, AgentKey: "test-key12333", GroupID: uint64(12223), Result: "bbb", Log: "log-xxxx"}
			tt.e.AddEvent(event1)
			time.Sleep(2 * time.Second)

			// hookTerm test
			event2 := &KlevrEvent{EventType: AgentDisconnect, AgentKey: "aaa-test-aaaakey12333", GroupID: uint64(912223), Result: "cc", Log: "log-xxxx"}
			tt.e.AddEvent(event2)
		})
	}

	time.Sleep(10 * time.Second)

	fmt.Println("end")
}
