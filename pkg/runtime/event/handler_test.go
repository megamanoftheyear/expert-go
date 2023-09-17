package event

import (
	"expert-go/pkg/interchange/metadata"
	"expert-go/pkg/unique"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestHandler_Emit(t *testing.T) {
	type fixture struct {
		testValue    string
		writeCounter atomic.Int32
	}

	onEventAction := func(f *fixture, value string) Action[unique.Key] {
		return func(m *metadata.Metadata[unique.Key]) {
			f.testValue = value
			f.writeCounter.Add(1)
		}
	}

	action := func(events ...unique.Key) func(handler *Handler[unique.Key]) {
		return func(handler *Handler[unique.Key]) {
			for _, event := range events {
				handler.Emit(event, nil)
			}
		}
	}

	const (
		event unique.Key = "event"
		group unique.Key = "group"
		chain unique.Key = "chain"
	)

	var (
		groupSubs = []unique.Key{"group_sub_1", "group_sub_2", "group_sub_3"}
		chainSubs = []unique.Key{"chain_sub_1", "chain_sub_2", "chain_sub_3"}
	)

	tests := []struct {
		name               string
		action             func(handler *Handler[unique.Key])
		handler            func(f *fixture) *Handler[unique.Key]
		expectedValue      string
		expectedWriteCount int32
	}{
		{
			name: "test on event action",
			handler: func(f *fixture) *Handler[unique.Key] {
				handler := NewHandler[unique.Key]()

				return handler.OnEvent(event, onEventAction(f, "event_test_value"))
			},
			action:             action(event),
			expectedValue:      "event_test_value",
			expectedWriteCount: 1,
		},
		{
			name: "test on group action",
			handler: func(f *fixture) *Handler[unique.Key] {
				handler := NewHandler[unique.Key]()
				g := NewGroup[unique.Key](group, groupSubs...)

				return handler.OnGroup(g, onEventAction(f, "group_test_value"))
			},
			action:             action(append(groupSubs, group)...),
			expectedValue:      "group_test_value",
			expectedWriteCount: int32(len(groupSubs) + 1),
		},
		{
			name: "test on chain action",
			handler: func(f *fixture) *Handler[unique.Key] {
				handler := NewHandler[unique.Key]()
				c, _ := NewChain[unique.Key](chain, chainSubs...)

				return handler.OnChain(c, onEventAction(f, "chain_test_value"))
			},
			action:             action(append(chainSubs, chain)...),
			expectedValue:      "chain_test_value",
			expectedWriteCount: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &fixture{}

			handler := test.handler(f)
			test.action(handler)
			assert.Equal(t, test.expectedValue, f.testValue)
			assert.Equal(t, test.expectedWriteCount, f.writeCounter.Load())
		})
	}
}
