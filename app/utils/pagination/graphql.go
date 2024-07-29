package pagination

import (
	"github.com/graphql-go/graphql"
)

const (
	DefaultLimit = 25
	MaxLimit     = 50
)

// PaginationType wraps a list type with pagination fields
func PaginationType[T any](object *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Paginated" + object.Name(),
		Fields: graphql.Fields{
			"pageInfo": {Type: PaginationInfo},
			"edges": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: object.Name() + "Edge",
					Fields: graphql.Fields{
						"cursor": &graphql.Field{
							Type: graphql.String,
						},
						"node": &graphql.Field{
							Type: graphql.NewNonNull(object),
							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								return p.Source.(Edge[T]).Node, nil
							},
						},
					},
				}))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*CursorPagination[T]).Edges, nil
				},
			},
		},
	})
}

// PaginationInputType defines the input type for pagination
var PaginationInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "PaginationInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"cursor": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"first": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"limit": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
		},
	},
)

var PagiantionArgs = &graphql.ArgumentConfig{
	Type: PaginationInputType,
}

var PaginationInfo = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"totalCount": &graphql.Field{
			Type: graphql.Int,
		},
		"hasNext": &graphql.Field{
			Type: graphql.Boolean,
		},
		"hasPrev": &graphql.Field{
			Type: graphql.Boolean,
		},
		"nextCursor": &graphql.Field{
			Type: graphql.String,
		},
		"prevCursor": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// PaginationOptionsFromArgs extracts pagination options from GraphQL arguments
func PaginationOptionsFromArgs[T any](args map[string]interface{}, model []T) Options[T] {
	limit := DefaultLimit // default limit
	cursor := ""

	if paginationArgs, ok := args["pagination"].(map[string]interface{}); ok {
		if l, ok := paginationArgs["limit"].(int); ok {
			if l < MaxLimit {
				limit = l
			} else {
				limit = MaxLimit
			}
		}

		if l, ok := paginationArgs["first"].(int); ok {
			if l < MaxLimit {
				limit = l
			} else {
				limit = MaxLimit
			}
		}

		if c, ok := paginationArgs["cursor"].(string); ok {
			cursor = c
		}
	}

	return Options[T]{
		Limit:  limit,
		Model:  model,
		Cursor: cursor,
	}
}
