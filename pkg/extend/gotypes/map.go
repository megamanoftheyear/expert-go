package gotypes

import (
	"expert-go/pkg/constraints"
	"github.com/alphadose/haxmap"
	"sync"
)

type Map[Key constraints.Hashable, Value any] struct{ container *haxmap.Map[Key, Value] }

func (m *Map[Key, Value]) Len() int                                { return int(m.container.Len()) }
func (m *Map[Key, Value]) Get(key Key) (Value, bool)               { return m.container.Get(key) }
func (m *Map[Key, Value]) Set(key Key, value Value)                { m.container.Set(key, value) }
func (m *Map[Key, Value]) Del(key Key)                             { m.container.Del(key) }
func (m *Map[Key, Value]) Each(fn func(key Key, value Value) bool) { m.container.ForEach(fn) }

type MapBuilder[Key constraints.Hashable, Value any] struct {
	lock *sync.Mutex

	container *Map[Key, Value]
}

func NewMapBuilder[Key constraints.Hashable, Value any]() *MapBuilder[Key, Value] {
	return &MapBuilder[Key, Value]{lock: &sync.Mutex{}}
}

func (builder *MapBuilder[Key, Value]) SetHaxmap(filledMap *haxmap.Map[Key, Value]) *MapBuilder[Key, Value] {
	builder.init()

	filledMap.ForEach(func(key Key, value Value) bool {
		builder.container.container.Set(key, value)

		return true
	})

	return builder
}

func (builder *MapBuilder[Key, Value]) SetMap(filledMap map[Key]Value) *MapBuilder[Key, Value] {
	builder.init()

	for key, value := range filledMap {
		builder.container.container.Set(key, value)
	}

	return builder
}

func (builder *MapBuilder[Key, Value]) SetSlice(slice []Key, setValue func(idx int, element Key) Value) *MapBuilder[Key, Value] {
	builder.init()

	for idx, key := range slice {
		builder.container.container.Set(key, setValue(idx, key))
	}

	return builder
}

func (builder *MapBuilder[Key, Value]) SetSliceOfDefaults(slice ...Key) *MapBuilder[Key, Value] {
	builder.init()

	for _, key := range slice {
		builder.container.container.Set(key, builder.defaultValue())
	}

	return builder
}

func (builder *MapBuilder[Key, Value]) defaultValue() Value { return *new(Value) }

func (builder *MapBuilder[Key, Value]) Build() *Map[Key, Value] {
	builder.init()

	builder.lock.Lock()
	defer builder.lock.Unlock()

	container := builder.container
	builder.container = nil

	return container
}

func (builder *MapBuilder[Key, Value]) init() {
	if builder.container != nil {
		return
	}

	builder.container = &Map[Key, Value]{
		container: haxmap.New[Key, Value](),
	}
}
