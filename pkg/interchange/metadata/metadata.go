package metadata

import (
	"expert-go/pkg/data-util/object"
	"expert-go/pkg/unique"
	"github.com/alphadose/haxmap"
	"reflect"
)

const allowedDepth = 15

type Schema map[unique.Key]reflect.Kind

type Metadata struct {
	metadata *haxmap.Map[unique.Key, any]
	schema   Schema
}

func NewMetadata(schema Schema) *Metadata {
	return &Metadata{
		schema:   schema,
		metadata: haxmap.New[unique.Key, any](),
	}
}

func (metadata *Metadata) Get(key unique.Key) (any, bool) { return metadata.metadata.Get(key) }
func (metadata *Metadata) Set(key unique.Key, value any) error {
	dump := object.GetDump(value)

	if dump.FieldsDepth() > allowedDepth {
		return &ObjectDepthError{
			typ:     dump.Type(),
			current: dump.FieldsDepth(),
			max:     allowedDepth,
		}
	}

	if len(metadata.schema) > 0 {
		t, exist := metadata.schema[key]
		if !exist {
			return &SchemaFieldError{
				schema:          metadata.schema,
				unexpectedField: key.String(),
			}
		}

		if reflect.TypeOf(value).Kind() != t {
			return &SchemaValueError{
				schema:   metadata.schema,
				value:    value,
				expected: t,
			}
		}
	}

	metadata.metadata.Set(key, value)

	return nil
}
