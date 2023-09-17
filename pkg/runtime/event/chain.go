package event

import (
	"expert-go/pkg/constraints"
	"expert-go/pkg/extend/gotypes"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"sync/atomic"
)

const (
	chainMaxLength = 20
)

type Chain[Key constraints.Hashable] struct {
	key Key

	events []Key

	registry *chainRegistry[Key]
}

func NewChain[Key constraints.Hashable](key Key, events ...Key) (*Chain[Key], error) {
	if eventsCount := len(events); eventsCount > chainMaxLength {
		return nil, NewChainMaxLengthError(eventsCount, eventsCount)
	}

	return &Chain[Key]{
		key:      key,
		events:   events,
		registry: newChainRegistry(events...),
	}, nil
}

func (chain *Chain[Key]) Key() Key      { return chain.key }
func (chain *Chain[Key]) Events() []Key { return chain.events }
func (chain *Chain[Key]) Include(event Key) bool {
	if !slices.Contains(chain.events, event) && chain.key != event {
		return false
	}

	registered, err := chain.Register(event)
	if err != nil {
		return false
	}

	return registered
}

func (chain *Chain[Key]) Register(event Key) (bool, error) {
	if chain.key == event {
		return true, nil
	}

	err := chain.registry.next(event)

	if IsChainSequenceError[Key](err) {
		return false, nil
	}

	if err != nil {
		return false, errors.WithMessage(err, "attempt to register next event")
	}

	return chain.registry.isComplete(), nil
}

type chainRegistry[Key constraints.Hashable] struct {
	registry *gotypes.Map[Key, int]

	idx atomic.Int32
}

func newChainRegistry[Key constraints.Hashable](events ...Key) *chainRegistry[Key] {
	setValue := func(idx int, _ Key) int { return idx }
	registry := gotypes.NewMapBuilder[Key, int]().SetSlice(events, setValue).Build()

	return &chainRegistry[Key]{registry: registry}
}

func (registry *chainRegistry[Key]) next(event Key) error {
	if registry.registry.Len() == 0 {
		return NewEndOfSequenceChainError(event)
	}

	expectedIdx, exist := registry.registry.Get(event)

	if !exist {
		return NewChainSequenceViolationError(registry.findExpectedKey(), event)
	}

	currentIdx := int(registry.idx.Load())
	if expectedIdx != currentIdx {
		return NewChainSequenceNumberViolationError(event, expectedIdx, currentIdx)
	}

	registry.registry.Del(event)
	registry.idx.Add(1)

	return nil
}

func (registry *chainRegistry[Key]) isComplete() bool { return registry.registry.Len() == 0 }

func (registry *chainRegistry[Key]) findExpectedKey() (k Key) {
	registry.registry.Each(func(key Key, value int) bool {
		if value == int(registry.idx.Load()) {
			k = key
		}

		return true
	})

	return k
}
