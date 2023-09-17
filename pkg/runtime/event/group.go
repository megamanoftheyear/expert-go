package event

import (
	"expert-go/pkg/constraints"
	"slices"
)

type Group[Key constraints.Hashable] struct {
	key Key

	events []Key
}

func NewGroup[Key constraints.Hashable](
	key Key,
	events ...Key,
) *Group[Key] {
	return &Group[Key]{
		key:    key,
		events: events,
	}
}

func (group *Group[Key]) Key() Key      { return group.key }
func (group *Group[Key]) Events() []Key { return group.events }
func (group *Group[Key]) Include(event Key) bool {
	return slices.Contains(group.events, event) || group.key == event
}
