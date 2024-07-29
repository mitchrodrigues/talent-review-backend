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
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/role"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/teams"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/pagination"
	"gorm.io/gorm"
)

func HandledCircularDeps(fnc func() *graphql.Object) *graphql.Object {
	return fnc()
}

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

	EmployeeRoleGQLType = graphql.NewObject(graphql.ObjectConfig{
		Name: "EmployeeRole",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(EmployeeRole).ID, nil
				},
			},
			"title": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(EmployeeRole).Title, nil
				},
			},
			"level": {
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(EmployeeRole).Level, nil
				},
			},
			"track": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					switch p.Source.(EmployeeRole).Track {
					case role.Manager:
						return "manager", nil
					default:
						return "ic", nil
					}
				},
			},
		},
	})

	EmployeeGQLType = graphql.NewObject(graphql.ObjectConfig{
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
			"role": {
				Type: EmployeeRoleGQLType,
				Resolve: gql.NewHandler(gql.Options{
					Handler: func(wctx golly.WebContext, params gql.Params) (interface{}, error) {
						emp := params.Source.(Employee)

						return wctx.Loader().Fetch(
							wctx.Context,
							fmt.Sprintf("employee_role:%s", emp.EmployeeRoleID), func(golly.Context) (interface{}, error) {
								if emp.Role.ID != uuid.Nil {
									return emp.Role, nil
								}

								return Service(wctx.Context).FindRoleByID(wctx.Context, emp.EmployeeRoleID)
							})
					},
				}),
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
			"leadID": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					team := p.Source.(Team)
					if team.LeadID == nil {
						return nil, nil
					}
					return team.LeadID.String(), nil
				},
			},
			"employees": {
				Type: graphql.NewNonNull(graphql.NewList(EmployeeGQLType)),
				Resolve: gql.NewHandler(gql.Options{
					Handler: func(wctx golly.WebContext, params gql.Params) (interface{}, error) {
						team := params.Source.(Team)

						return Service(wctx.Context).FindEmployeesForTeam(wctx.Context, team.ID)
					},
				}),
			},
			"lead": {
				Type: EmployeeGQLType,
				Resolve: gql.NewHandler(gql.Options{
					Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
						if leadID := params.Source.(Team).LeadID; leadID != nil && *leadID != uuid.Nil {
							return ctx.Loader().
								Fetch(ctx.Context, fmt.Sprintf("employee:%s", leadID.String()),
									func(golly.Context) (interface{}, error) {
										return Service(ctx.Context).FindEmployeeByID(ctx.Context, *leadID)
									})
						}
						return nil, nil
					},
				}),
			},
		},
	})

	query = graphql.Fields{
		//********** Roles ***************//
		"employeeRoles": &graphql.Field{
			Name: "employeeRoles",
			Args: graphql.FieldConfigArgument{
				"pagination": pagination.PagiantionArgs,
				"filter":     roleFilter.Args,
			},
			Type: pagination.PaginationType[EmployeeRole](EmployeeRoleGQLType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					scopes, err := roleFilter.Scopes(ctx.Context, params.Args["filter"])
					if err != nil {
						return nil, err
					}

					return pagination.
						NewCursorPaginationFromArgs(
							params.Args,
							[]EmployeeRole{},
							scopes...,
						).
						SetScopes(common.OrganizationIDScopeForContext(ctx.Context, "employee_roles")).
						Paginate(ctx.Context)
				},
			}),
		},

		//********** Employees ***************//
		"employees": &graphql.Field{
			Name: "employees",
			Args: graphql.FieldConfigArgument{
				"pagination": pagination.PagiantionArgs,
				"filter":     employeeFilter.Args,
			},
			Type: pagination.PaginationType[Employee](EmployeeGQLType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					scopes, err := employeeFilter.Scopes(ctx.Context, params.Args["filter"])
					if err != nil {
						return nil, err
					}

					records, err := pagination.
						NewCursorPaginationFromArgs(
							params.Args,
							[]Employee{},
							scopes...,
						).
						SetScopes(common.OrganizationIDScopeForContext(ctx.Context, "employees")).
						SetScopes(func(db *gorm.DB) *gorm.DB {
							return db.Preload("Role").Preload("Team")
						}).
						Paginate(ctx.Context)

					return records.Cache(ctx.Context, "employee:%s"), err
				},
			}),
		},
		"employee": &graphql.Field{
			Name: "employee",
			Args: graphql.FieldConfigArgument{"id": {Type: graphql.NewNonNull(graphql.String)}},
			Type: EmployeeGQLType,
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					employee, err := Service(ctx.Context).FindEmployeeByID(ctx.Context, id)
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
			Type: graphql.NewList(EmployeeGQLType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					ident := identity.FromContext(ctx.Context)

					scopes, err := employeeFilter.Scopes(ctx.Context, params.Args["filter"])
					if err != nil {
						return nil, err
					}

					return Service(ctx.Context).FindEmployeesByManagerUserID(
						ctx.Context,
						ident.UID,
						scopes...)
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
					scopes, err := teamFilter.Scopes(ctx.Context, params.Args["filter"])
					if err != nil {
						return nil, err
					}

					return pagination.
						NewCursorPaginationFromArgs(params.Args, []Team{}, scopes...).
						SetScopes(common.OrganizationIDScopeForContext(ctx.Context)).
						Paginate(ctx.Context)
				},
			}),
		},
	}

	createTeamInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateTeamInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":   {Type: graphql.NewNonNull(graphql.String)},
			"leadID": {Type: graphql.String},
		},
	})

	updateTeamInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateTeamInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":   {Type: graphql.NewNonNull(graphql.String)},
			"leadID": {Type: graphql.String},
		},
	})

	createEmployeeInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateEmployeeInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":       {Type: graphql.NewNonNull(graphql.String)},
			"email":      {Type: graphql.NewNonNull(graphql.String)},
			"workerType": {Type: graphql.NewNonNull(workerType)},
			"teamID":     {Type: graphql.String},
			"roleID":     {Type: graphql.String},
			"managerID":  {Type: graphql.String},
		},
	})

	updateEmployeeInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateEmployeeInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":       {Type: graphql.String},
			"email":      {Type: graphql.String},
			"teamID":     {Type: graphql.String},
			"workerType": {Type: workerType},
			"managerID":  {Type: graphql.String},
			"roleID":     {Type: graphql.String},
		},
	})

	createEmployeeRoleInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateEmployeeRoleInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"title": {Type: graphql.NewNonNull(graphql.String)},
			"track": {Type: graphql.NewNonNull(graphql.String)},
			"level": {Type: graphql.NewNonNull(graphql.Int)},
		},
	})

	updateEmployeeRoleInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateEmployeeRoleInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"title": {Type: graphql.String},
			"track": {Type: graphql.String},
			"level": {Type: graphql.Int},
		},
	})

	mutations = graphql.Fields{
		//********** EmployeeRoles ***************//
		"createEmployeeRole": &graphql.Field{
			Type: EmployeeRoleGQLType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: createEmployeeRoleInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					var empRole EmployeeRole

					err := eventsource.Call(ctx.Context, &empRole.Aggregate, role.Create{
						Title:          params.Input["title"].(string),
						Level:          int(params.Input["level"].(int)),
						Track:          role.EmployeeType(params.Input["track"].(string)),
						OrganizationID: identity.FromContext(ctx.Context).OrganizationID,
					}, params.Metadata())

					if err != nil {
						return nil, err
					}

					return empRole, nil
				},
			}),
		},

		"updateEmployeeRole": &graphql.Field{
			Type: EmployeeRoleGQLType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: updateEmployeeRoleInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {

					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					empRole, err := Service(ctx.Context).FindRoleByID(ctx.Context, id)
					if err != nil {
						return nil, errors.WrapNotFound(err)
					}

					track, _ := helpers.ExtractArg[string](params.Input, "track")
					level, _ := helpers.ExtractArg[int](params.Input, "level")
					title, _ := helpers.ExtractArg[string](params.Input, "title")

					err = eventsource.Call(ctx.Context, &empRole.Aggregate, role.Update{
						Title: title,
						Level: level,
						Track: track,
					}, params.Metadata())

					if err != nil {
						return nil, err
					}

					return empRole, nil
				},
			}),
		},

		//********** Employees ***************//
		"createEmployee": &graphql.Field{
			Name: "createEmployee",
			Type: EmployeeGQLType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: createEmployeeInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					ident := identity.FromContext(ctx.Context)

					var emp Employee
					var teamID uuid.UUID

					if val, err := helpers.ExtractAndParseUUID(params.Input, "teamID"); err == nil {
						teamID = val
					}

					managerID, _ := helpers.ExtractAndParseUUID(params.Input, "managerID")
					roleID, _ := helpers.ExtractAndParseUUID(params.Input, "roleID")

					err := eventsource.Call(ctx.Context, &emp.Aggregate, employee.Create{
						Name:           params.Input["name"].(string),
						Email:          params.Input["email"].(string),
						TeamID:         teamID,
						EmployeeRoleID: roleID,
						WorkerType:     employee.EmployeeWorkerType(params.Input["workerType"].(string)),
						OrganizationID: ident.OrganizationID,
						ManagerID:      managerID,
					}, params.Metadata())

					return emp, err
				},
			}),
		},

		"updateEmployee": &graphql.Field{
			Type: EmployeeGQLType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: updateEmployeeInputType},
			},
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					var teamID uuid.UUID
					var managerID uuid.UUID
					var roleID uuid.UUID

					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					emp, err := Service(ctx.Context).FindEmployeeByID(ctx.Context, id)
					if err != nil {
						return nil, err
					}

					if val, err := helpers.ExtractAndParseUUID(params.Input, "teamID"); err == nil {
						teamID = val
					}

					if val, err := helpers.ExtractAndParseUUID(params.Input, "managerID"); err == nil {
						managerID = val
					}

					if val, err := helpers.ExtractAndParseUUID(params.Input, "roleID"); err == nil {
						roleID = val
					}

					err = eventsource.Call(ctx.Context, &emp.Aggregate, employee.Update{
						Name:           params.Input["name"].(string),
						Email:          params.Input["email"].(string),
						TeamID:         teamID,
						ManagerID:      managerID,
						EmployeeRoleID: roleID,
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
					var leadID *uuid.UUID

					if l, err := helpers.ExtractAndParseUUID(params.Input, "leadID"); err == nil {
						leadID = &l
					}

					err := eventsource.Call(ctx.Context, &team.Aggregate, teams.Create{
						Name:           params.Input["name"].(string),
						LeadID:         leadID,
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

					var leadID *uuid.UUID

					if l, err := helpers.ExtractAndParseUUID(params.Input, "leadID"); err == nil {
						leadID = &l
					}

					name, _ := helpers.ExtractArg[string](params.Input, "name")

					var team Team

					err = orm.DB(ctx.Context).Model(team).
						Scopes(common.OrganizationIDScopeForContext(ctx.Context)).
						Find(&team, "id = ?", id).
						Error

					if err != nil {
						return nil, err
					}

					err = eventsource.Call(ctx.Context, &team.Aggregate, teams.Update{
						Name:   name,
						LeadID: leadID,
					}, params.Metadata())

					return team, err
				},
			}),
		},
	}
)

// TODO Refactor this into chunks where we can easily define these duplications
func AddCircularDependencies() {
	/**** Employee *****/
	EmployeeGQLType.AddFieldConfig("team", &graphql.Field{
		Type: teamGQLType,
		Resolve: gql.NewHandler(gql.Options{
			Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
				if emp := params.Source.(Employee); emp.TeamID != nil {
					if *emp.TeamID != uuid.Nil {
						return golly.LoadData(ctx.Context, fmt.Sprintf("team:%s", emp.TeamID.String()), func(golly.Context) (Team, error) {
							if emp.Team.ID != uuid.Nil {
								return emp.Team, nil
							}

							return Service(ctx.Context).FindTeamByID(ctx.Context, *emp.TeamID)
						})
					}
				}
				return nil, nil
			},
		}),
	})

	EmployeeGQLType.AddFieldConfig("manager", &graphql.Field{
		Type: EmployeeGQLType,
		Resolve: gql.NewHandler(gql.Options{
			Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
				if managerID := params.Source.(Employee).ManagerID; managerID != nil {
					return golly.LoadData(ctx.Context, fmt.Sprintf("employee:%s", managerID.String()), func(golly.Context) (Employee, error) {
						return Service(ctx.Context).FindEmployeeByID(ctx.Context, *managerID)
					})
				}
				return nil, nil
			},
		}),
	})

}

func InitGraphQL() {
	AddCircularDependencies()

	gql.RegisterQuery(query)
	gql.RegisterMutation(mutations)
}
