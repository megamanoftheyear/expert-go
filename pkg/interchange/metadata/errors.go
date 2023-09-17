package metadata

import (
	"expert-go/pkg/constraints"
	"expert-go/pkg/data-util/object"
	"expert-go/pkg/unique"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type ObjectDepthError struct {
	typ          object.Type
	current, max int
}

func (err *ObjectDepthError) Typ() object.Type { return err.typ }
func (err *ObjectDepthError) Current() int     { return err.current }
func (err *ObjectDepthError) Max() int         { return err.max }

func (err *ObjectDepthError) Is(target error) bool {
	return errors.As(target, &ObjectDepthError{})
}

func (err *ObjectDepthError) Error() string {
	return fmt.Sprintf(
		"invalid object `%s` depth `%d`, max `%d`",
		err.typ.String(),
		err.current,
		err.max,
	)
}

type SchemaError interface{ Schema() Schema[unique.Key] }

type SchemaFieldError[Key constraints.Hashable] struct {
	schema          *Schema[Key]
	unexpectedField Key
}

func (err *SchemaFieldError[Key]) Schema() *Schema[Key] { return err.schema }
func (err *SchemaFieldError[Key]) UnexpectedField() Key { return err.unexpectedField }

func (err *SchemaFieldError[Key]) Is(target error) bool {
	return errors.As(target, &SchemaFieldError[Key]{})
}

func (err *SchemaFieldError[Key]) Error() string {
	return fmt.Sprintf("unexpected schema field `%s`", err.unexpectedField)
}

type SchemaValueError[Key constraints.Hashable] struct {
	schema   *Schema[Key]
	value    any
	expected reflect.Kind
}

func (err *SchemaValueError[Key]) Schema() *Schema[Key]   { return err.schema }
func (err *SchemaValueError[Key]) Value() any             { return err.value }
func (err *SchemaValueError[Key]) Expected() reflect.Kind { return err.expected }

func (err *SchemaValueError[Key]) Is(target error) bool {
	return errors.As(target, &SchemaValueError[Key]{})
}

func (err *SchemaValueError[Key]) Error() string {
	return fmt.Sprintf(
		"unexpected schema value `%v`(%T), expected type `%s`",
		err.value,
		err.value,
		err.expected,
	)
}
