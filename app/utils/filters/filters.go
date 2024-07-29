package filters

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/golly-go/golly"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FieldType represents the type of a field in GraphQL and the corresponding database field name.
type FieldType struct {
	Kind           reflect.Kind
	DBFieldName    string
	BoolExpression string
	Clauses        func(gctx golly.Context) []clause.Expression
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
	// Define the field enum
	fieldEnumValues := graphql.EnumValueConfigMap{}
	for field := range fields {
		fieldEnumValues[field] = &graphql.EnumValueConfig{
			Value: field,
		}
	}

	fieldEnum := graphql.NewEnum(graphql.EnumConfig{
		Name:   fmt.Sprintf("%sFieldEnum", name),
		Values: fieldEnumValues,
	})

	// Define the filter input object
	filterInputObject := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: fmt.Sprintf("%sFilterInput", name),
		Fields: graphql.InputObjectConfigFieldMap{
			"field": &graphql.InputObjectFieldConfig{
				Type: fieldEnum,
			},
			"operator": &graphql.InputObjectFieldConfig{
				Type: operatorType,
			},
			"value": &graphql.InputObjectFieldConfig{
				Type: graphql.String, // Placeholder, will be dynamically validated
			},
		},
	})

	return &graphql.ArgumentConfig{
		Type: graphql.NewList(filterInputObject),
	}
}

// Scopes generates GORM scopes from GraphQL arguments
func (f *Filter) Scopes(context golly.Context, filterArgs interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
	scopes := []func(*gorm.DB) *gorm.DB{}
	joins := map[string]bool{}

	args, ok := filterArgs.([]interface{})
	if !ok {
		return scopes, nil
	}

	for _, arg := range args {
		argMap, ok := arg.(map[string]interface{})
		if !ok {
			continue
		}

		field, ok := argMap["field"].(string)
		if !ok {
			continue
		}

		operator, ok := argMap["operator"].(string)
		if !ok {
			operator = "is" // Default to "is" operator
		}

		fieldType, ok := f.Fields[field]
		if !ok {
			return scopes, fmt.Errorf("invalild field %s", field)
		}

		var value interface{}

		rawValue := argMap["value"]
		if operator != "has" && operator != "has not" {
			if rawValue == nil {
				return scopes, fmt.Errorf("invalild value supplied %s", rawValue)
			}

			v, err := convertValue(fieldType.Kind, rawValue)
			if err != nil {
				return scopes, err
			}

			value = v
		}

		if fieldType.Clauses != nil {
			fieldClauses := fieldType.Clauses(context)

			golly.Each(fieldClauses, func(fClause clause.Expression) {
				switch clse := fClause.(type) {
				case clause.Join:
					if exists := joins[clse.Table.Alias]; exists {
						return
					}
					joins[clse.Table.Alias] = true
					scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
						return db.Joins(clauseToQuery(clse, db))
					})
				case clause.Where:
					scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
						return db.Where(clauseToQuery(clse, db))
					})
				default:
					context.Logger().Debugf("Invalid clause type %#v\n", clse)
					return
				}
			})
		}

		if fieldType.BoolExpression != "" {
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("(%s) = ?", fieldType.BoolExpression), value)
			})
			continue
		}

		sqlOperator, val := extractOperatorAndValue(operator, value)

		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s %s ?", fieldType.DBFieldName, sqlOperator), val)
		})
	}

	return scopes, nil
}

var operatorMap = map[string]string{
	"has":            "IS NOT NULL",
	"has not":        "IS NULL",
	"contains":       "ILIKE",
	"contains_exact": "LIKE",
	"starts_with":    "ILIKE",
	"ends_with":      "ILIKE",
	"is":             "=",
	"is not":         "<>",
	"greater than":   ">",
	"less than":      "<",
}

// Define the GraphQL enum using the operatorMap
var operatorType = graphql.NewEnum(graphql.EnumConfig{
	Name: "FilterOperator",
	Values: graphql.EnumValueConfigMap{
		"HAS": &graphql.EnumValueConfig{
			Value: "has",
		},
		"HAS_NOT": &graphql.EnumValueConfig{
			Value: "has not",
		},
		"CONTAINS": &graphql.EnumValueConfig{
			Value: "contains",
		},
		"CONTAINS_EXACT": &graphql.EnumValueConfig{
			Value: "contains_exact",
		},
		"STARTS_WITH": &graphql.EnumValueConfig{
			Value: "starts_with",
		},
		"ENDS_WITH": &graphql.EnumValueConfig{
			Value: "ends_with",
		},
		"IS": &graphql.EnumValueConfig{
			Value: "is",
		},
		"IS_NOT": &graphql.EnumValueConfig{
			Value: "is not",
		},
		"GREATER_THAN": &graphql.EnumValueConfig{
			Value: "greater than",
		},
		"LESS_THAN": &graphql.EnumValueConfig{
			Value: "less than",
		},
	},
})

// convertValue converts a value to the correct type based on reflect.Kind
func convertValue(kind reflect.Kind, value interface{}) (interface{}, error) {
	strVal, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("value is not a string")
	}

	switch kind {
	case reflect.String:
		return strVal, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.ParseInt(strVal, 10, 64)
	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(strVal, 64)
	case reflect.Bool:
		return strconv.ParseBool(strVal)
	default:
		return nil, fmt.Errorf("unsupported kind: %s", kind)
	}
}

// extractOperatorAndValue extracts the operator and value from the provided argument
func extractOperatorAndValue(operator string, value interface{}) (string, interface{}) {
	sqlOperator, exists := operatorMap[operator]
	if !exists {
		sqlOperator = "="
	}

	ret := value
	switch sqlOperator {
	case "ILIKE", "LIKE":
		if operator == "contains" {
			if strVal, ok := value.(string); ok {
				ret = fmt.Sprintf("%%%s%%", strVal)
			}
		} else if operator == "starts_with" {
			if strVal, ok := value.(string); ok {
				ret = fmt.Sprintf("%s%%", strVal)
			}
		} else if operator == "ends_with" {
			if strVal, ok := value.(string); ok {
				ret = fmt.Sprintf("%%%s", strVal)
			}
		}
	case "IS NOT NULL", "IS NULL":
		ret = nil
	}

	return sqlOperator, ret
}

func clauseToQuery(clse clause.Expression, db *gorm.DB) string {
	stmt := &gorm.Statement{SQL: strings.Builder{}, DB: db}

	clse.Build(stmt)

	return stmt.SQL.String()
}
