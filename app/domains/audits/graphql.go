package audits

import (
	"fmt"
	"strings"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/gql"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/pagination"
)

var (
	auditUserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "AuditUser",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(accounts.User).ID, nil
				},
			},
			"name": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := p.Source.(accounts.User)

					return fmt.Sprintf("%s %s", user.FirstName, user.LastName), nil
				},
			},
		},
	})

	auditGQLType = graphql.NewObject(graphql.ObjectConfig{
		Name: "AuditLog",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Event).ID, nil
				},
			},
			"event": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return strings.Split(p.Source.(Event).Type, ".")[1], nil
				},
			},
			"objectType": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return strings.Split(p.Source.(Event).AggregateType, ".")[0], nil
				},
			},
			"objectID": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Event).ID, nil
				},
			},
			"changes": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					b, err := p.Source.(Event).RawData.MarshalJSON()
					return string(b), err
				},
			},
			"eventAt": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Event).CreatedAt, nil
				},
			},
			"user": {
				Type: auditUserType,
				Resolve: gql.NewHandler(gql.Options{
					Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
						event := params.Source.(Event)

						if event.UserID == nil {
							return nil, nil
						}

						user, err := accounts.FindUserByID(ctx.Context, event.UserID.String())
						if err != nil || user.ID != uuid.Nil {
							return nil, nil
						}

						return user, nil
					},
				}),
			},
		},
	})

	query = graphql.Fields{
		"auditLogs": {
			Name: "auditLogs",
			Type: pagination.PaginationType[Event](auditGQLType),
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: pagination.PaginationInputType,
				},
			},

			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					return pagination.
						NewCursorPaginationFromArgs(
							params.Args,
							[]Event{},
							common.OrganizationIDScopeForContext(ctx.Context),
						).
						Paginate(ctx.Context)
				},
			}),
		},
	}

	mutations = graphql.Fields{}
)

func InitGraphQL() {
	gql.RegisterQuery(query)
	gql.RegisterMutation(mutations)
}
