package employees

import (
	"fmt"
	"reflect"

	"github.com/golly-go/golly"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/role"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/filters"
	"gorm.io/gorm/clause"
)

// I will start by saying i probaly hate this as its super verbose but
// it works and is flexible meh
/************** Role Filter **************/

var (
	roleFilter = filters.NewFilter("EmployeeRole", map[string]filters.FieldType{
		"title": {Kind: reflect.String, DBFieldName: "employee_roles.title"},
	})
)

/************** Employee Filter **************/
var (
	employeeFilter = filters.NewFilter("Employee", map[string]filters.FieldType{
		"name":   {Kind: reflect.String, DBFieldName: "employees.name"},
		"teamID": {Kind: reflect.String, DBFieldName: "employees.team_id"},
		"isManager": {
			Kind:           reflect.Bool,
			BoolExpression: fmt.Sprintf("role.track = '%s'", role.Manager),
			Clauses: func(golly.Context) []clause.Expression {
				return []clause.Expression{
					clause.Join{
						Table: clause.Table{Name: EmployeeRole{}.TableName(), Alias: "role"},
						ON: clause.Where{
							Exprs: []clause.Expression{
								clause.Eq{
									Column: "role.id",
									Value:  clause.Column{Table: "employees", Name: "employee_role_id"},
								},
								clause.Eq{
									Column: "role.organization_id",
									Value:  clause.Column{Table: "employees", Name: "organization_id"},
								},
							},
						},
					},
				}
			},
		},
		"manager": {
			Kind:        reflect.String,
			DBFieldName: "manager.name",
			Clauses: func(golly.Context) []clause.Expression {
				return []clause.Expression{
					clause.Join{
						Table: clause.Table{Name: Employee{}.TableName(), Alias: "manager"},
						ON: clause.Where{
							Exprs: []clause.Expression{
								clause.Eq{
									Column: "manager.id",
									Value:  clause.Column{Table: clause.CurrentTable, Name: "manager_id"},
								},
								clause.Eq{
									Column: "manager.organization_id",
									Value:  clause.Column{Table: clause.CurrentTable, Name: "organization_id"},
								},
							},
						},
					},
				}
			},
		},
		"team": {
			Kind:        reflect.String,
			DBFieldName: "team.name",
			Clauses: func(ctx golly.Context) []clause.Expression {
				return []clause.Expression{
					clause.Join{
						Table: clause.Table{Name: Team{}.TableName(), Alias: "team"},
						ON: clause.Where{
							Exprs: []clause.Expression{
								clause.Eq{
									Column: "team.id",
									Value:  clause.Column{Table: clause.CurrentTable, Name: "team_id"},
								},
								clause.Eq{
									Column: "team.organization_id",
									Value:  clause.Column{Table: clause.CurrentTable, Name: "organization_id"},
								},
							},
						},
					},
				}
			},
		},
	})
)

/************** Team Filter **************/
var (
	teamFilter = filters.NewFilter("Team", map[string]filters.FieldType{
		"name": {Kind: reflect.String, DBFieldName: "name"},
		"lead": {
			Kind:        reflect.String,
			DBFieldName: "lead.name",
			Clauses: func(golly.Context) []clause.Expression {
				return []clause.Expression{
					clause.Join{
						Table: clause.Table{Name: "employees", Alias: "lead"},
						ON: clause.Where{
							Exprs: []clause.Expression{
								clause.Eq{Column: "lead.id", Value: clause.PrimaryColumn},
								clause.Eq{
									Column: "lead.organization_id",
									Value:  clause.Column{Table: "employees", Name: "organization_id"},
								},
							},
						},
					},
				}
			},
		},
	})
)
