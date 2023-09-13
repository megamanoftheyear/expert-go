package metadata

import (
	"expert-go/pkg/data-util/object"
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

type SchemaError interface{ Schema() Schema }

type SchemaFieldError struct {
	schema          Schema
	unexpectedField string
}

func (err *SchemaFieldError) Schema() Schema          { return err.schema }
func (err *SchemaFieldError) UnexpectedField() string { return err.unexpectedField }

func (err *SchemaFieldError) Is(target error) bool {
	return errors.As(target, &SchemaFieldError{})
}

func (err *SchemaFieldError) Error() string {
	return fmt.Sprintf("unexpected schema field `%s`", err.unexpectedField)
}

type SchemaValueError struct {
	schema   Schema
	value    any
	expected reflect.Kind
}

func (err *SchemaValueError) Schema() Schema         { return err.schema }
func (err *SchemaValueError) Value() any             { return err.value }
func (err *SchemaValueError) Expected() reflect.Kind { return err.expected }

func (err *SchemaValueError) Is(target error) bool {
	return errors.As(target, &SchemaValueError{})
}

func (err *SchemaValueError) Error() string {
	return fmt.Sprintf(
		"unexpected schema value `%v`(%T), expected type `%s`",
		err.value,
		err.value,
		err.expected,
	)
}
