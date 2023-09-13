package metadata

import (
	"expert-go/pkg/data-util/object"
	"expert-go/pkg/unique"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMetadata_Set(t *testing.T) {
	tests := []struct {
		name          string
		key           unique.Key
		value         any
		schema        Schema
		expectedError error
	}{
		{
			name:  "invalid depth struct",
			key:   "key",
			value: invalidDepthStruct{},
			expectedError: &ObjectDepthError{
				typ:     object.Struct,
				current: 16,
				max:     15,
			},
		},
		{
			name:   "invalid key of schema",
			key:    "invalid_key",
			value:  10,
			schema: Schema{"valid_key": reflect.Bool},
			expectedError: &SchemaFieldError{
				schema:          Schema{"valid_key": reflect.Bool},
				unexpectedField: "invalid_key",
			},
		},
		{
			name:   "invalid value of schema",
			key:    "key",
			value:  true,
			schema: Schema{"key": reflect.String},
			expectedError: &SchemaValueError{
				schema:   Schema{"key": reflect.String},
				value:    true,
				expected: reflect.String,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metadata := NewMetadata(test.schema)

			err := metadata.Set(test.key, test.value)
			if test.expectedError != nil {
				t.Logf(test.expectedError.Error())
				assert.EqualError(t, err, test.expectedError.Error())
			}
		})
	}
}

type invalidDepthStruct struct { //1
	Key struct { //2
		Key struct { //3
			Key struct { //4
				Key struct { //5
					Key struct { //6
						Key struct { //7
							Key struct { //8
								Key struct { //9
									Key struct { //10
										Key struct { //11
											Key struct { //12
												Key struct { //13
													Key struct { //14
														Key struct { //15
															Key struct { //16
																Key int
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
