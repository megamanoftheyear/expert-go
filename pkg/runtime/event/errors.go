package event

import (
	"expert-go/pkg/constraints"
	"fmt"
	"github.com/pkg/errors"
)

type ChainMaxLengthError struct{ current, max int }

func NewChainMaxLengthError(current int, max int) *ChainMaxLengthError {
	return &ChainMaxLengthError{current: current, max: max}
}

func (err *ChainMaxLengthError) Current() int { return err.current }
func (err *ChainMaxLengthError) Max() int     { return err.max }

func (err *ChainMaxLengthError) Error() string {
	msg := "maximum number[%d] of event chain elements[%d] exceeded"

	return fmt.Sprintf(msg, err.max, err.current)
}

type ChainSequenceError[Key constraints.Hashable] interface{ IsChainSequenceError() bool }

func IsChainSequenceError[Key constraints.Hashable](target error) bool {
	var err ChainSequenceError[Key]

	return errors.As(target, &err)
}

type ChainSequenceViolationError[Key constraints.Hashable] struct{ expected, actual Key }

func NewChainSequenceViolationError[Key constraints.Hashable](expected Key, actual Key) *ChainSequenceViolationError[Key] {
	return &ChainSequenceViolationError[Key]{expected: expected, actual: actual}
}

func (err *ChainSequenceViolationError[Key]) Expected() Key              { return err.expected }
func (err *ChainSequenceViolationError[Key]) Actual() Key                { return err.actual }
func (err *ChainSequenceViolationError[Key]) IsChainSequenceError() bool { return true }

func (err *ChainSequenceViolationError[Key]) Error() string {
	msg := "the sequence chain was broken: expected `%s`, actual `%s`"

	return fmt.Sprintf(msg, err.expected, err.actual)
}

type ChainSequenceNumberViolationError[Key constraints.Hashable] struct {
	event            Key
	expected, actual int
}

func NewChainSequenceNumberViolationError[Key constraints.Hashable](event Key, expected, actual int) *ChainSequenceNumberViolationError[Key] {
	return &ChainSequenceNumberViolationError[Key]{event: event, expected: expected, actual: actual}
}

func (err *ChainSequenceNumberViolationError[Key]) Event() Key                 { return err.event }
func (err *ChainSequenceNumberViolationError[Key]) Expected() int              { return err.expected }
func (err *ChainSequenceNumberViolationError[Key]) Actual() int                { return err.actual }
func (err *ChainSequenceNumberViolationError[Key]) IsChainSequenceError() bool { return true }

func (err *ChainSequenceNumberViolationError[Key]) Error() string {
	msg := "the event `%s` has violated its expected[%d] number on the current[%d]"

	return fmt.Sprintf(msg, err.event, err.expected, err.actual)
}

type EndOfSequenceChainError[Key constraints.Hashable] struct{ event Key }

func NewEndOfSequenceChainError[Key constraints.Hashable](event Key) *EndOfSequenceChainError[Key] {
	return &EndOfSequenceChainError[Key]{event: event}
}

func (err *EndOfSequenceChainError[Key]) Event() Key                 { return err.event }
func (err *EndOfSequenceChainError[Key]) IsChainSequenceError() bool { return true }

func (err *EndOfSequenceChainError[Key]) Error() string {
	return fmt.Sprintf("attempt to write an event `%v` to a closed circuit", err.event)
}
