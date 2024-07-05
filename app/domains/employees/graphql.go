package employees

import (
	"fmt"
	"strings"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/gql"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/teams"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/filters"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/pagination"
	"gorm.io/gorm"
)

// TODO: Break these up into smaller files but for now put this here cause its quicker to dev against

var (
	workerType = graphql.NewEnum(graphql.EnumConfig{
		Name: "WorkerType",
		Values: graphql.EnumValueConfigMap{
			"agency": {Value: "AC"},
			"direct": {Value: "DC"},
			"fte":    {Value: "FTE"},
		},
	})

	employeeTrack = graphql.NewEnum(graphql.EnumConfig{
		Name: "EmployeeType",
		Values: graphql.EnumValueConfigMap{
			"ic":      {Value: "IC"},
			"manager": {Value: "MNG"},
		},
	})

	EmployeeSimplified = graphql.NewObject(graphql.ObjectConfig{
		Name: "EmployeeInfo",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Employee).ID, nil
				},
			},
			"name": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Employee).Name, nil
				},
			},
		},
	})

	employeeGQLType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Employee",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Employee).ID, nil
				},
			},
			"name": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Employee).Name, nil
				},
			},
			"email": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Employee).Email, nil
				},
			},
			"level": {
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Employee).Level, nil
				},
			},
			"type": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					switch p.Source.(Employee).Type {
					case employee.Manager:
						return "manager", nil
					default:
						return "ic", nil
					}
				},
			},
			"title": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Employee).Title, nil
				},
			},
			"workerType": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					switch strings.TrimSpace(string(p.Source.(Employee).WorkerType)) {
					case string(employee.AgencyContractor):
						return "agency", nil
					case string(employee.DirectContractor):
						return "direct", nil
					default:
						return "fte", nil
					}

				},
			},
			"team": {
				Type: teamGQLType,
				Resolve: gql.NewHandler(gql.Options{
					Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
						if teamID := params.Source.(Employee).TeamID; teamID != nil {
							return golly.LoadData(ctx.Context, fmt.Sprintf("team:%s", teamID), func(golly.Context) (Team, error) {
								return FindTeamByID(ctx.Context, *teamID)
							})
						}
						return nil, nil
					},
				}),
			},
		},
	})

	teamGQLType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Team",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Team).ID, nil
				},
			},
			"name": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Team).Name, nil
				},
			},
			"managerID": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					team := p.Source.(Team)

					if team.ManagerID == uuid.Nil {
						return nil, nil
					}
					return team.ManagerID.String(), nil
				},
			},
			"manager": {
				Type: EmployeeSimplified,
				Resolve: gql.NewHandler(gql.Options{
					Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
						team := params.Source.(Team)
						var manager Employee

						if team.ManagerID == uuid.Nil {
							return nil, nil
						}

						err := orm.
							NewDB(ctx.Context).
							Model(manager).
							Find(&manager, "id = ?", team.ManagerID).
							Error

						if err != nil {
							return nil, err
						}

						return manager, nil
					},
				}),
			},
		},
	})

	employeeFilter = filters.NewFilter("Employee", map[string]filters.FieldType{
		"manager": {GraphQLType: employeeTrack, DBFieldName: "type"},
		"name":    {GraphQLType: graphql.String, DBFieldName: "name", Wildcard: true},
	})

	teamFilter = filters.NewFilter("Team", map[string]filters.FieldType{
		"name": {GraphQLType: graphql.String, DBFieldName: "name", Wildcard: true},
	})

	query = graphql.Fields{
		//********** Employees ***************//
		"employees": &graphql.Field{
			Name: "employees",
			Args: graphql.FieldConfigArgument{
				"pagination": pagination.PagiantionArgs,
				"filter":     employeeFilter.Args,
			},
			Type: pagination.PaginationType[Employee](employeeGQLType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					return pagination.
						NewCursorPaginationFromArgs(
							params.Args,
							[]Employee{},
							employeeFilter.Scopes(params.Args["filter"])...,
						).
						SetScopes(common.OrganizationIDScopeForContext(ctx.Context)).
						SetScopes(func(db *gorm.DB) *gorm.DB {
							return db.Preload("Team.Manager")
						}).
						Paginate(ctx.Context)
				},
			}),
		},
		"employee": &graphql.Field{
			Name: "employee",
			Args: graphql.FieldConfigArgument{"id": {Type: graphql.NewNonNull(graphql.String)}},
			Type: employeeGQLType,
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					employee, err := FindEmployeeByID(ctx.Context, id)
					if err != nil {
						return nil, err
					}

					return employee, nil
				},
			}),
		},

		"subordinates": &graphql.Field{
			Name: "subordinates",
			Args: graphql.FieldConfigArgument{
				"filter": employeeFilter.Args,
			},
			Type: graphql.NewList(employeeGQLType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					ident := identity.FromContext(ctx.Context)

					return FindEmployeesByManagersUserID(
						ctx.Context,
						ident.UID,
						employeeFilter.Scopes(params.Args["filter"])...)
				},
			}),
		},

		//********** TEAMS ***************//

		"teams": &graphql.Field{
			Name: "teams",
			Args: graphql.FieldConfigArgument{
				"pagination": pagination.PagiantionArgs,
				"filter":     teamFilter.Args,
			},
			Type: pagination.PaginationType[Team](teamGQLType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					return pagination.
						NewCursorPaginationFromArgs(
							params.Args,
							[]Team{},
							teamFilter.Scopes(params.Args["filter"])...).
						SetScopes(common.OrganizationIDScopeForContext(ctx.Context)).
						Paginate(ctx.Context)
				},
			}),
		},
	}

	createTeamInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateTeamInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":      {Type: graphql.NewNonNull(graphql.String)},
			"managerID": {Type: graphql.String},
		},
	})

	updateTeamInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateTeamInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":      {Type: graphql.NewNonNull(graphql.String)},
			"managerID": {Type: graphql.String},
		},
	})

	createEmployeeInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateEmployeeInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":       {Type: graphql.NewNonNull(graphql.String)},
			"email":      {Type: graphql.NewNonNull(graphql.String)},
			"workerType": {Type: graphql.NewNonNull(workerType)},
			"title":      {Type: graphql.String},
			"teamID":     {Type: graphql.String},
			"level":      {Type: graphql.NewNonNull(graphql.Int)},
			"manager":    {Type: graphql.NewNonNull(graphql.Boolean)},
		},
	})

	updateEmployeeInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateEmployeeInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":       {Type: graphql.NewNonNull(graphql.String)},
			"email":      {Type: graphql.NewNonNull(graphql.String)},
			"title":      {Type: graphql.String},
			"teamID":     {Type: graphql.String},
			"level":      {Type: graphql.NewNonNull(graphql.Int)},
			"workerType": {Type: graphql.NewNonNull(workerType)},
		},
	})

	mutations = graphql.Fields{
		//********** Employees ***************//
		"createEmployee": &graphql.Field{
			Name: "createEmployee",
			Type: employeeGQLType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: createEmployeeInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					ident := identity.FromContext(ctx.Context)

					var emp Employee
					var teamID uuid.UUID
					var title string

					if val, err := helpers.ExtractAndParseUUID(params.Input, "teamID"); err == nil {
						teamID = val
					}

					if val, err := helpers.ExtractArg[string](params.Input, "title"); err == nil {
						title = val
					}

					err := eventsource.Call(ctx.Context, &emp.Aggregate, employee.Create{
						Name:           params.Input["name"].(string),
						Email:          params.Input["email"].(string),
						Manager:        params.Input["manager"].(bool),
						Level:          params.Input["level"].(int),
						TeamID:         teamID,
						Title:          title,
						OrganizationID: ident.OrganizationID,
						WorkerType:     employee.EmployeeWorkerType(params.Input["workerType"].(string)),
					}, params.Metadata())

					return emp, err
				},
			}),
		},

		"updateEmployee": &graphql.Field{
			Type: employeeGQLType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: updateEmployeeInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					var teamID uuid.UUID
					var title string

					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					emp, err := FindEmployeeByID(ctx.Context, id)
					if err != nil {
						return nil, err
					}

					if val, err := helpers.ExtractAndParseUUID(params.Input, "teamID"); err == nil {
						teamID = val
					}

					if val, err := helpers.ExtractArg[string](params.Input, "title"); err == nil {
						title = val
					}

					err = eventsource.Call(ctx.Context, &emp.Aggregate, employee.Update{
						Name:       params.Input["name"].(string),
						Email:      params.Input["email"].(string),
						Level:      params.Input["level"].(int),
						WorkerType: employee.EmployeeWorkerType(params.Input["workerType"].(string)),
						TeamID:     teamID,
						Title:      title,
					}, params.Metadata())

					return emp, err

				},
			}),
		},

		//********** TEAMS ***************//
		"createTeam": &graphql.Field{
			Name: "createTeam",
			Type: teamGQLType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: createTeamInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					ident := identity.FromContext(ctx.Context)

					var team Team
					var managerID uuid.UUID

					managerID, err := helpers.ExtractAndParseUUID(params.Input, "managerID")
					if err != nil {
						return nil, err
					}

					err = eventsource.Call(ctx.Context, &team.Aggregate, teams.CreateTeam{
						Name:           params.Input["name"].(string),
						ManagerID:      managerID,
						OrganizationID: ident.OrganizationID,
					}, params.Metadata())

					return team, err
				},
			}),
		},

		"updateTeam": &graphql.Field{
			Name: "createTeam",
			Type: teamGQLType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: updateTeamInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, errors.WrapNotFound(err)
					}

					managerID, _ := helpers.ExtractAndParseUUID(params.Input, "managerID")
					name, _ := helpers.ExtractArg[string](params.Input, "name")

					var team Team

					err = orm.DB(ctx.Context).Model(team).
						Scopes(common.OrganizationIDScopeForContext(ctx.Context)).
						Find(&team, "id = ?", id).
						Error

					if err != nil {
						return nil, err
					}

					err = eventsource.Call(ctx.Context, &team.Aggregate, teams.UpdateTeam{
						Name:      name,
						ManagerID: managerID,
					}, params.Metadata())

					return team, err
				},
			}),
		},
	}
)

func InitGraphQL() {
	gql.RegisterQuery(query)
	gql.RegisterMutation(mutations)
}
