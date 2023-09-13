package object

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetDump(t *testing.T) {
	type m map[string]any

	tests := []struct {
		name          string
		obj           any
		expectedType  Type
		expectedDepth int
	}{
		{
			name:          "map",
			obj:           m{"1": m{"2": 1, "3": m{"4": m{"5": m{"6": 2}}}}},
			expectedType:  Map,
			expectedDepth: 4,
		},
		{
			name: "struct",
			obj: struct {
				Key1 struct {
					Key2 struct{ Key3 struct{ Key4 int } }
				}
			}{},
			expectedType:  Struct,
			expectedDepth: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dump := GetDump(test.obj)

			require.NotNil(t, dump)
			assert.Equal(t, test.expectedType.String(), dump.typ.String())
			assert.Equal(t, test.expectedDepth, dump.fieldsDepth)
		})
	}
}
