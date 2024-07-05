package filters

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

// FieldType represents the type of a field in GraphQL and the corresponding database field name.
type FieldType struct {
	GraphQLType graphql.Input
	DBFieldName string
	Wildcard    bool
}

// Filter struct to hold field mappings and their types
type Filter struct {
	Name   string
	Fields map[string]FieldType
	Args   *graphql.ArgumentConfig
}

// NewFilter initializes a new Filter
func NewFilter(name string, fields map[string]FieldType) *Filter {
	return &Filter{
		Name:   name,
		Fields: fields,
		Args:   graphQLArgs(name, fields),
	}
}

// graphQLArgs generates a GraphQL argument config from the Filter fields
func graphQLArgs(name string, fields map[string]FieldType) *graphql.ArgumentConfig {
	graphqlFields := graphql.InputObjectConfigFieldMap{}

	for field, fieldType := range fields {
		graphqlFields[field] = &graphql.InputObjectFieldConfig{Type: fieldType.GraphQLType}
	}

	inputObject := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   fmt.Sprintf("%sFilter", name),
		Fields: graphqlFields,
	})

	return &graphql.ArgumentConfig{
		Type: inputObject,
	}
}

// Scopes generates GORM scopes from GraphQL arguments
func (f *Filter) Scopes(filterArgs interface{}) (scopes []func(*gorm.DB) *gorm.DB) {
	args, ok := filterArgs.(map[string]interface{})
	if !ok {
		return []func(*gorm.DB) *gorm.DB{}
	}

	for arg, value := range args {
		fieldType, ok := f.Fields[arg]
		if !ok {
			continue
		}

		if strValue, ok := value.(string); ok && fieldType.Wildcard {
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("LOWER(%s) LIKE ?", fieldType.DBFieldName), strings.ToLower(strValue)+"%")
			})
			continue
		}

		val := value
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", fieldType.DBFieldName), val)
		})
	}

	return scopes
}
