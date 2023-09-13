package object

import (
	"github.com/google/go-cmp/cmp"
	"github.com/nanorobocop/go-recursive"
	"golang.org/x/exp/slices"
	"reflect"
)

func GetDump(v any) *Dump {
	dump := &Dump{typ: getType(v)}

	levelFields := make(map[int]struct{})

	var lastObject any

	recursive.Go(&v, func(value any, level int) interface{} {
		if cmp.Equal(lastObject, value) || level == 0 {
			return value
		}

		if _, exists := levelFields[level]; !exists && getType(value).WithFields() {
			dump.fieldsDepth++
		}

		levelFields[level] = struct{}{}
		lastObject = value

		return value
	})

	return dump
}

func getType(v any) Type {
	t := reflect.TypeOf(v)

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if typ, exists := reflectTypes[t.Kind()]; exists {
		return typ
	}

	return Primitive
}

type Type uint

func (typ Type) WithFields() bool { return slices.Contains([]Type{Map, Struct}, typ) }
func (typ Type) IsValid() bool    { return typeLower < typ && typ < typeUpper }
func (typ Type) String() string {
	if !typ.IsValid() {
		return "unknown"
	}

	return types[typ]
}

const (
	typeLower Type = iota

	Struct
	Slice
	Map
	Primitive
	Interface

	typeUpper
)

var (
	types = map[Type]string{
		Struct:    "Struct",
		Slice:     "Slice",
		Map:       "Map",
		Primitive: "Primitive",
		Interface: "Interface",
	}

	reflectTypes = map[reflect.Kind]Type{
		reflect.Interface: Interface,
		reflect.Struct:    Struct,
		reflect.Slice:     Slice,
		reflect.Map:       Map,
	}
)

type Dump struct {
	typ Type

	fieldsDepth int
}

func (dump *Dump) Type() Type       { return dump.typ }
func (dump *Dump) FieldsDepth() int { return dump.fieldsDepth }
