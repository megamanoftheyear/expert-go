package event

import (
	"expert-go/pkg/constraints"
	"expert-go/pkg/interchange/metadata"
)

type Event[Key constraints.Hashable] struct {
	key Key

	metadata *metadata.Metadata[Key]
}

func NewEvent[Key constraints.Hashable](
	key Key,
	metadata *metadata.Metadata[Key],
) *Event[Key] {
	return &Event[Key]{
		key:      key,
		metadata: metadata,
	}
}

func (event *Event[Key]) Key() Key                          { return event.key }
func (event *Event[Key]) Metadata() *metadata.Metadata[Key] { return event.metadata }

type MultiEvent[Key constraints.Hashable] interface {
	Include(event Key) bool
	Key() Key
}
