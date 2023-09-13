package unique

import (
	"github.com/alphadose/haxmap"
	"sync"
)

type Map[Value any] struct{ container *haxmap.Map[Key, Value] }

func (m Map[Value]) Len() int { return int(m.container.Len()) }

type MapBuilder[Value any] struct {
	lock *sync.Mutex

	container *Map[Value]
}

func NewMapBuilder[Value any]() *MapBuilder[Value] {
	return &MapBuilder[Value]{lock: &sync.Mutex{}}
}

func (builder *MapBuilder[Value]) WithFilledHaxmap(filledMap *haxmap.Map[Key, Value]) *MapBuilder[Value] {
	builder.init()

	filledMap.ForEach(func(key Key, value Value) bool {
		builder.container.container.Set(key, value)

		return true
	})

	return builder
}

func (builder *MapBuilder[Value]) WithFilledMap(filledMap map[Key]Value) *MapBuilder[Value] {
	builder.init()

	for key, value := range filledMap {
		builder.container.container.Set(key, value)
	}

	return builder
}

func (builder *MapBuilder[Value]) Build() *Map[Value] {
	builder.init()

	builder.lock.Lock()
	defer builder.lock.Unlock()

	container := builder.container
	builder.container = nil

	return container
}

func (builder *MapBuilder[Value]) init() {
	if builder.container != nil {
		return
	}

	builder.container = &Map[Value]{
		container: haxmap.New[Key, Value](),
	}
}
