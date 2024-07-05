package employees

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"gorm.io/gorm"
)

func FindEmployeeByUserID(gctx golly.Context, userID uuid.UUID) (Employee, error) {
	var employees Employee

	err := orm.DB(gctx).
		Model(Employee{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&employees, "user_id = ?", userID).
		Error

	return employees, err

}

func FindEmployeesByIDS(gctx golly.Context, ids uuid.UUIDs) ([]Employee, error) {
	var employees []Employee

	err := orm.DB(gctx).
		Model(Employee{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&employees, "id IN ?", ids).
		Error

	return employees, err

}

func FindEmployeesForTeam(gctx golly.Context, teamID uuid.UUID, excludeEmployees ...uuid.UUID) ([]Employee, error) {

	var employees []Employee

	err := orm.DB(gctx).
		Model(Employee{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Where("id NOT IN ?", excludeEmployees).
		Find(&employees, "team_id = ?", teamID).
		Error

	return employees, err

}

func FindEmployeesByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error) {
	var teams []Team

	err := orm.DB(gctx).
		Model(Team{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Scopes(scopes...).
		Preload("Employees").
		Find(&teams, "manager_id = ?", managerID).
		Error

	return golly.Flatten(
		golly.Map(teams, func(team Team) []Employee {
			return team.Employees
		}),
	), err
}

func FindEmployeeIDsByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]uuid.UUID, error) {
	ident := identity.FromContext(gctx)

	var employeeIDs []uuid.UUID

	err := orm.DB(gctx).
		Model(Employee{}).
		Scopes(scopes...).
		Joins("JOIN teams ON employees.team_id = teams.id").
		Where("employees.organization_id = ?", ident.OrganizationID).
		Where("teams.manager_id = ?", managerID).
		Pluck("employees.id", &employeeIDs).
		Error

	return employeeIDs, err
}

func FindEmployeeIDsByManagerUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]uuid.UUID, error) {
	myRecord, err := FindEmployeeByUserID(gctx, userID)
	if err != nil || myRecord.ID == uuid.Nil {
		return nil, err
	}

	return FindEmployeeIDsByManagerID(gctx, myRecord.ID, scopes...)
}

func FindEmployeesByManagersUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error) {
	myRecord, err := FindEmployeeByUserID(gctx, userID)
	if err != nil || myRecord.ID == uuid.Nil {
		return nil, err
	}

	return FindEmployeesByManagerID(
		gctx,
		myRecord.ID,
		scopes...)

}

func FindEmployeeByEmailAndOrganizationID(gctx golly.Context, email string, organizationID uuid.UUID) (Employee, error) {
	var emp Employee

	err := orm.DB(gctx).
		Model(emp).
		Scopes(common.OrganizationIDScope(organizationID)).
		Where("email = ?", email).
		Find(&emp).
		Error

	return emp, errors.WrapGeneric(err)

}

func FindEmployeeByID(gctx golly.Context, id uuid.UUID) (Employee, error) {
	var emp Employee

	err := orm.
		DB(gctx).
		Model(emp).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&emp, "id = ?", id).
		Error

	return emp, errors.WrapNotFound(err)
}

func FindEmployeeByID_Unsafe(gctx golly.Context, id uuid.UUID) (Employee, error) {
	var emp Employee

	err := orm.
		DB(gctx).
		Model(emp).
		Find(&emp, "id = ?", id).
		Error

	return emp, errors.WrapNotFound(err)
}

func FindTeamsByOrganizationID(gctx golly.Context, organizationID uuid.UUID) ([]Team, error) {
	var teams []Team

	err := orm.
		DB(gctx).
		Model(Team{}).
		Scopes(common.OrganizationIDScope(organizationID)).
		Scopes(DefaultPreloads).
		Find(&teams).
		Error

	return teams, err
}

func FindTeamByID(gctx golly.Context, id uuid.UUID) (Team, error) {
	var team Team

	err := orm.
		DB(gctx).
		Model(Team{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&team, "id = ?", id).
		Error

	return team, err
}

func DefaultPreloads(db *gorm.DB) *gorm.DB {
	return db.Preload("Manager")
}
