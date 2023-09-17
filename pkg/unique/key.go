package unique

import (
	"context"
	"sync/atomic"
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
}

func NewSharedKey(readCloser bool) *SharedKey {
	return &SharedKey{readCloser: readCloser}
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
	if !key.closed.Load() && key.readCloser {
		key.closed.Store(true)
	}

	key.Each(func(shard Key) { k += shard })

	return k
}

func (key *SharedKey) Each(fn func(Key)) *SharedKey {
	if !key.closed.Load() && key.readCloser {
		key.closed.Store(true)
	}

	for _, shard := range key.key {
		fn(shard)
	}

	return key
}

func (key *SharedKey) Exist(shard Key) (exist bool) {
	if !key.closed.Load() && key.readCloser {
		key.closed.Store(true)
	}

	key.Each(func(key Key) {
		if key == shard {
			exist = true
		}
	})

	return exist
}

func (key *SharedKey) Close() *SharedKey {
	if !key.closed.Load() && key.readCloser {
		key.closed.Store(true)
	}

	return key
}
