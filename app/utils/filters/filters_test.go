package filters

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

// TestGraphQLArgs tests the GraphQLArgs method of the Filter struct
func TestGraphQLArgs(t *testing.T) {
	tests := []struct {
		name     string
		filter   *Filter
		expected *graphql.ArgumentConfig
	}{
		{
			name: "single field",
			filter: NewFilter("Test", map[string]FieldType{
				"field1": {GraphQLType: graphql.String},
			}),
			expected: &graphql.ArgumentConfig{
				Type: graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "TestFilter",
					Fields: graphql.InputObjectConfigFieldMap{
						"field1": &graphql.InputObjectFieldConfig{Type: graphql.String},
					},
				}),
			},
		},
		{
			name: "multiple fields",
			filter: NewFilter("Test", map[string]FieldType{
				"field1": {GraphQLType: graphql.String},
				"field2": {GraphQLType: graphql.Int},
			}),
			expected: &graphql.ArgumentConfig{
				Type: graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "TestFilter",
					Fields: graphql.InputObjectConfigFieldMap{
						"field1": &graphql.InputObjectFieldConfig{Type: graphql.String},
						"field2": &graphql.InputObjectFieldConfig{Type: graphql.Int},
					},
				}),
			},
		},
		{
			name:   "no fields",
			filter: NewFilter("Empty", map[string]FieldType{}),
			expected: &graphql.ArgumentConfig{
				Type: graphql.NewInputObject(graphql.InputObjectConfig{
					Name:   "EmptyFilter",
					Fields: graphql.InputObjectConfigFieldMap{},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.Args

			assert.Equal(t, tt.expected.Type.Name(), result.Type.Name())

			expectedFields := tt.expected.Type.(*graphql.InputObject).Fields()
			resultFields := result.Type.(*graphql.InputObject).Fields()

			for fieldName, expectedField := range expectedFields {
				resultField, exists := resultFields[fieldName]
				assert.True(t, exists)
				assert.Equal(t, expectedField.Type, resultField.Type)
			}
		})
	}
}
