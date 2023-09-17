package metadata

import (
	"expert-go/pkg/constraints"
	"expert-go/pkg/data-util/object"
	"expert-go/pkg/extend/gotypes"
	"github.com/alphadose/haxmap"
	"reflect"
)

const allowedDepth = 15

type Schema[Key constraints.Hashable] struct {
	m *gotypes.Map[Key, reflect.Kind]
}

func NewSchema[Key constraints.Hashable](builder *gotypes.MapBuilder[Key, reflect.Kind]) *Schema[Key] {
	if builder == nil {
		builder = gotypes.NewMapBuilder[Key, reflect.Kind]()
	}

	return &Schema[Key]{m: builder.Build()}
}

func (schema *Schema[Key]) Get(key Key) (reflect.Kind, bool) { return schema.m.Get(key) }
func (schema *Schema[Key]) Len() int                         { return schema.m.Len() }
func (schema *Schema[Key]) Set(key Key, kind reflect.Kind) *Schema[Key] {
	schema.m.Set(key, kind)

	return schema
}

type Metadata[Key constraints.Hashable] struct {
	metadata *haxmap.Map[Key, any]
	schema   *Schema[Key]
}

func NewMetadata[Key constraints.Hashable](schema *Schema[Key]) *Metadata[Key] {
	if schema == nil {
		schema = NewSchema[Key](nil)
	}

	return &Metadata[Key]{
		schema:   schema,
		metadata: haxmap.New[Key, any](),
	}
}

func (metadata *Metadata[Key]) Get(key Key) (any, bool) { return metadata.metadata.Get(key) }
func (metadata *Metadata[Key]) Set(key Key, value any) error {
	dump := object.GetDump(value)

	if dump.FieldsDepth() > allowedDepth {
		return &ObjectDepthError{
			typ:     dump.Type(),
			current: dump.FieldsDepth(),
			max:     allowedDepth,
		}
	}

	if metadata.schema.Len() > 0 {
		t, exist := metadata.schema.Get(key)
		if !exist {
			return &SchemaFieldError[Key]{
				schema:          metadata.schema,
				unexpectedField: key,
			}
		}

		if reflect.TypeOf(value).Kind() != t {
			return &SchemaValueError[Key]{
				schema:   metadata.schema,
				value:    value,
				expected: t,
			}
		}
	}

	metadata.metadata.Set(key, value)

	return nil
}
