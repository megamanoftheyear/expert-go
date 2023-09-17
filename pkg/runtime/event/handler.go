package event

import (
	"expert-go/pkg/constraints"
	"expert-go/pkg/data-util/cmp"
	"expert-go/pkg/extend/gotypes"
	"expert-go/pkg/interchange/metadata"
)

type Action[Key constraints.Hashable] func(mtd *metadata.Metadata[Key])

func (action Action[Key]) IsActive() bool { return action != nil }

type Handler[Key constraints.Hashable] struct {
	on     *gotypes.Map[Key, Action[Key]]
	groups []*Group[Key]
	chains []*Chain[Key]
}

func NewHandler[Key constraints.Hashable]() *Handler[Key] {
	return &Handler[Key]{on: gotypes.NewMapBuilder[Key, Action[Key]]().Build()}
}

func (handler *Handler[Key]) OnEvent(event Key, action Action[Key]) *Handler[Key] {
	handler.on.Set(event, action)

	return handler
}

func (handler *Handler[Key]) OnGroup(group *Group[Key], action Action[Key]) *Handler[Key] {
	handler.OnEvent(group.Key(), action)

	return handler.RegisterGroup(group)
}

func (handler *Handler[Key]) OnChain(chain *Chain[Key], action Action[Key]) *Handler[Key] {
	handler.OnEvent(chain.Key(), action)

	return handler.RegisterChain(chain)
}

func (handler *Handler[Key]) RegisterGroup(group *Group[Key]) *Handler[Key] {
	handler.groups = append(handler.groups, group)

	return handler
}

func (handler *Handler[Key]) RegisterChain(chain *Chain[Key]) *Handler[Key] {
	handler.chains = append(handler.chains, chain)

	return handler
}

func (handler *Handler[Key]) Emit(key Key, mtd *metadata.Metadata[Key]) *Handler[Key] {
	if mtd == nil {
		mtd = metadata.NewMetadata[Key](nil)
	}

	handler.execEventActions(key, mtd)

	return handler
}

func (handler *Handler[Key]) execEventActions(key Key, mtd *metadata.Metadata[Key]) {
	if action, exist := handler.on.Get(key); exist {
		action(mtd)

		return
	}

	groupMultiEvents := gotypes.SliceOf[*Group[Key], MultiEvent[Key]](
		handler.groups, func(group *Group[Key]) MultiEvent[Key] { return group },
	)
	chainMultiEvents := gotypes.SliceOf[*Chain[Key], MultiEvent[Key]](
		handler.chains, func(chain *Chain[Key]) MultiEvent[Key] { return chain },
	)

	groups := handler.findMultiEvents(key, groupMultiEvents...)
	chains := handler.findMultiEvents(key, chainMultiEvents...)

	for _, multiEvent := range append(groups, chains...) {
		if cmp.Filled[Key](&multiEvent) {
			continue
		}

		if action, exist := handler.on.Get(multiEvent); exist {
			action(mtd)
		}
	}

}

func (handler *Handler[Key]) findMultiEvents(key Key, events ...MultiEvent[Key]) (multiEvents []Key) {
	for _, event := range events {
		if event.Include(key) {
			multiEvents = append(multiEvents, event.Key())
		}
	}

	return multiEvents
}
