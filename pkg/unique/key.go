package unique

import (
	"context"
	"expert-go/pkg/interchange/metadata"
	"expert-go/pkg/runtime/event"
	"sync/atomic"
)

const (
	EventKey   event.Key = "Key"
	EventEach  event.Key = "Each"
	EventExist event.Key = "Exist"
	EventClose event.Key = "Close"
)

type Key string

func (key Key) Empty() bool    { return key.Len() == 0 }
func (key Key) String() string { return string(key) }
func (key Key) Len() int       { return len(key) }

type Keygen interface{ Next(context.Context) Key }

type SharedKey struct {
	key []Key

	closed     atomic.Bool
	readCloser bool

	events *event.Handler
}

func NewSharedKey(readCloser bool) *SharedKey {
	key := &SharedKey{readCloser: readCloser}

	return key.initEventsHandler()
}

func (key *SharedKey) String() string { return key.Key().String() }

func (key *SharedKey) Add(shard Key) *SharedKey {
	if key.closed.Load() || key.Exist(shard) {
		return key
	}

	key.key = append(key.key, shard)

	return key
}

func (key *SharedKey) Key() (k Key) {
	key.events.Emit(EventKey, nil)
	key.Each(func(shard Key) { k += shard })

	return k
}

func (key *SharedKey) Each(fn func(Key)) *SharedKey {
	key.events.Emit(EventEach, nil)

	for _, shard := range key.key {
		fn(shard)
	}

	return key
}

func (key *SharedKey) Exist(shard Key) (exist bool) {
	key.events.Emit(EventExist, nil)

	key.Each(func(key Key) {
		if key == shard {
			exist = true
		}
	})

	return exist
}

func (key *SharedKey) Close() *SharedKey {
	key.events.Emit(EventClose, nil)

	return key
}

func (key *SharedKey) initEventsHandler() *SharedKey {
	handler := event.NewHandler()

	closeGroup := event.NewGroup(EventClose, EventKey, EventEach, EventExist)
	handler.OnGroup(closeGroup, func(*metadata.Metadata) {
		if !key.closed.Load() && key.readCloser {
			key.closed.Store(true)
		}
	})

	key.events = handler

	return key
}
