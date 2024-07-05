package reviews

import (
	"github.com/graphql-go/graphql"
)

var (
	query = graphql.Fields{
		// "users": {
		// 	Name: "users",
		// 	Args: graphql.FieldConfigArgument{
		// 		"pagination": &graphql.ArgumentConfig{
		// 			Type: pagination.PaginationInputType,
		// 		},
		// 	},
		// 	Type: pagination.PaginationType[User](userType),
		// 	Resolve: gql.NewHandler(gql.Options{
		// 		Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
		// 			return pagination.
		// 				NewCursorPaginationFromArgs(
		// 					params.Args,
		// 					[]User{},
		// 					common.OrganizationIDScopeForContext(ctx.Context)).
		// 				Paginate(ctx.Context)
		// 		},
		// 	}),
		// },
	}
)
