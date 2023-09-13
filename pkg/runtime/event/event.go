package event

import (
	"expert-go/pkg/interchange/metadata"
	"expert-go/pkg/unique"
	"github.com/alphadose/haxmap"
)

type Key unique.Key

func (key Key) Empty() bool { return unique.Key(key).Empty() }

type Group struct {
	key Key

	events []Key
}

func NewGroup(
	key Key,
	events ...Key,
) Group {
	return Group{
		key:    key,
		events: events,
	}
}

func (group *Group) Key() Key      { return group.key }
func (group *Group) Events() []Key { return group.events }

type Event struct {
	key Key

	metadata *metadata.Metadata
}

func NewEvent(
	key Key,
	metadata *metadata.Metadata,
) *Event {
	return &Event{
		key:      key,
		metadata: metadata,
	}
}

func (event *Event) Key() Key                     { return event.key }
func (event *Event) Metadata() *metadata.Metadata { return event.metadata }

type Action func(mtd *metadata.Metadata)

func (action Action) IsActive() bool { return action != nil }

type Handler struct {
	on     *haxmap.Map[Key, Action]
	groups []Group
}

func NewHandler() *Handler { return &Handler{} }

func (handler *Handler) OnEvent(event Key, action Action) *Handler {
	handler.on.Set(event, action)

	return handler
}

func (handler *Handler) OnGroup(group Group, action Action) *Handler {
	handler.OnEvent(group.Key(), action)

	return handler.RegisterGroup(group.Key(), group.Events()...)
}

func (handler *Handler) RegisterGroup(group Key, events ...Key) *Handler {
	handler.groups = append(handler.groups, Group{key: group, events: events})

	return handler
}

func (handler *Handler) Emit(key Key, mtd *metadata.Metadata) *Handler {
	if mtd == nil {
		mtd = metadata.NewMetadata(nil)
	}

	eventAction := handler.findEventAction(key)

	if eventAction.IsActive() {
		eventAction(mtd)
	}

	return handler
}

func (handler *Handler) findEventAction(key Key) Action {
	if action, exist := handler.on.Get(key); exist {
		return action
	}

	group := handler.findGroup(key)

	if group.Empty() {
		return nil
	}

	if action, exist := handler.on.Get(group); exist {
		return action
	}

	return nil
}

func (handler *Handler) findGroup(key Key) Key {
	for _, group := range handler.groups {
		for _, event := range group.Events() {

			if event == key {
				return group.Key()
			}

		}
	}

	return ""
}
