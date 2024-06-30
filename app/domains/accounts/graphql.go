package accounts

import (
	"fmt"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/gql"
	"github.com/graphql-go/graphql"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/pagination"
)

var (
	organizationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Organization",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Organization).ID, nil
				},
			},
			"name": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Organization).Name, nil
				},
			},
		},
	})

	userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(User).ID, nil
				},
			},
			"firstName": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(User).FirstName, nil
				},
			},
			"lastName": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(User).LastName, nil
				},
			},
			"email": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(User).Email, nil
				},
			},
			"invitedAt": {
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := p.Source.(User)
					if user.InvitedAt == nil {
						return nil, nil
					}
					return user.InvitedAt, nil
				},
			},
			"language": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "en-US", nil
				},
			},
			"organization": {
				Type: organizationType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(User).Organization, nil
				},
			},
		},
	})

	query = graphql.Fields{
		"users": {
			Name: "users",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: pagination.PaginationInputType,
				},
			},
			Type: pagination.PaginationType[User](userType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					return pagination.
						NewCursorPaginationFromArgs(params.Args, []User{}).
						Paginate(ctx.Context)
				},
			}),
		},
		"organization": {
			Name: "organization",
			Type: organizationType,
			Args: graphql.FieldConfigArgument{
				"id": {Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					ident := identity.FromContext(ctx.Context)
					if ident.OrganizationID != id {
						return nil, errors.WrapNotFound(fmt.Errorf("record not found"))
					}

					organization, err := FindOrganizationByID(ctx.Context, id)
					if err != nil {
						return nil, err
					}

					return organization, nil
				},
			}),
		},
		"me": {
			Name: "me",
			Type: userType,
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					user, err := FindUserByID(ctx.Context, params.Identity.UserID(), DefaultUserPreloads)

					if err != nil {
						return nil, err
					}

					return user, nil
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
