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

const (
	serviceCtxKey golly.ContextKeyT = "employeeService"
)

type EmployeeService interface {
	FindEmployeeByUserID(gctx golly.Context, userID uuid.UUID) (Employee, error)
	FindEmployeesByIDS(gctx golly.Context, ids uuid.UUIDs) ([]Employee, error)
	FindEmployeesForTeam(gctx golly.Context, teamID uuid.UUID, excludeEmployees ...uuid.UUID) ([]Employee, error)
	FindEmployeesByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error)
	FindEmployeeIDsByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) (uuid.UUIDs, error)
	FindEmployeeIDsByManagerUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) (uuid.UUIDs, error)
	FindEmployeesByManagersUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error)
	FindEmployeeByEmailAndOrganizationID(gctx golly.Context, email string, organizationID uuid.UUID) (Employee, error)
	FindEmployeeByID(gctx golly.Context, id uuid.UUID) (Employee, error)
	FindEmployeeEmailsBySearch(gctx golly.Context, name string) ([]string, error)
	FindEmployeeByID_Unsafe(gctx golly.Context, id uuid.UUID) (Employee, error)
	FindTeamsByOrganizationID(gctx golly.Context, organizationID uuid.UUID) ([]Team, error)
	FindTeamByID(gctx golly.Context, id uuid.UUID) (Team, error)
}

type DefaultEmployeeService struct{}

func (s DefaultEmployeeService) FindEmployeeByUserID(gctx golly.Context, userID uuid.UUID) (Employee, error) {
	var employees Employee

	err := orm.DB(gctx).
		Model(Employee{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&employees, "user_id = ?", userID).
		Error

	return employees, err

}

func (s DefaultEmployeeService) FindEmployeesByIDS(gctx golly.Context, ids uuid.UUIDs) ([]Employee, error) {
	var employees []Employee

	err := orm.DB(gctx).
		Model(Employee{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&employees, "id IN ?", ids).
		Error

	return employees, err

}

func (s DefaultEmployeeService) FindEmployeesForTeam(gctx golly.Context, teamID uuid.UUID, excludeEmployees ...uuid.UUID) ([]Employee, error) {

	var employees []Employee

	err := orm.DB(gctx).
		Model(Employee{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Where("id NOT IN ?", excludeEmployees).
		Find(&employees, "team_id = ?", teamID).
		Error

	return employees, err

}

func (s DefaultEmployeeService) FindEmployeesByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error) {
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

func (s DefaultEmployeeService) FindEmployeeIDsByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) (uuid.UUIDs, error) {
	ident := identity.FromContext(gctx)

	var employeeIDs uuid.UUIDs

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

func (s DefaultEmployeeService) FindEmployeeIDsByManagerUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) (uuid.UUIDs, error) {
	myRecord, err := s.FindEmployeeByUserID(gctx, userID)
	if err != nil || myRecord.ID == uuid.Nil {
		return nil, err
	}

	return s.FindEmployeeIDsByManagerID(gctx, myRecord.ID, scopes...)
}

func (s DefaultEmployeeService) FindEmployeesByManagersUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error) {
	myRecord, err := s.FindEmployeeByUserID(gctx, userID)
	if err != nil || myRecord.ID == uuid.Nil {
		return nil, err
	}

	return s.FindEmployeesByManagerID(
		gctx,
		myRecord.ID,
		scopes...)

}

func (s DefaultEmployeeService) FindEmployeeByEmailAndOrganizationID(gctx golly.Context, email string, organizationID uuid.UUID) (Employee, error) {
	var emp Employee

	err := orm.DB(gctx).
		Model(emp).
		Scopes(common.OrganizationIDScope(organizationID)).
		Where("email = ?", email).
		Find(&emp).
		Error

	return emp, errors.WrapGeneric(err)

}

func (s DefaultEmployeeService) FindEmployeeByID(gctx golly.Context, id uuid.UUID) (Employee, error) {
	var emp Employee

	err := orm.
		DB(gctx).
		Model(emp).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&emp, "id = ?", id).
		Error

	return emp, errors.WrapNotFound(err)
}

func (s DefaultEmployeeService) FindEmployeeEmailsBySearch(gctx golly.Context, name string) ([]string, error) {
	var emails []string

	err := orm.
		DB(gctx).
		Model(&Employee{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Select("DISTINCT(email) AS email").
		Where("LOWER(email) LIKE ?", name+"%").
		Pluck("email", &emails).
		Error

	return emails, err
}

func (s DefaultEmployeeService) FindEmployeeByID_Unsafe(gctx golly.Context, id uuid.UUID) (Employee, error) {
	var emp Employee

	err := orm.
		DB(gctx).
		Model(emp).
		Find(&emp, "id = ?", id).
		Error

	return emp, errors.WrapNotFound(err)
}

func (s DefaultEmployeeService) FindTeamsByOrganizationID(gctx golly.Context, organizationID uuid.UUID) ([]Team, error) {
	var teams []Team

	err := orm.
		DB(gctx).
		Model(Team{}).
		Scopes(common.OrganizationIDScope(organizationID)).
		Scopes(s.DefaultPreloads).
		Find(&teams).
		Error

	return teams, err
}

func (s DefaultEmployeeService) FindTeamByID(gctx golly.Context, id uuid.UUID) (Team, error) {
	var team Team

	err := orm.
		DB(gctx).
		Model(Team{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&team, "id = ?", id).
		Error

	return team, err
}

func (s DefaultEmployeeService) DefaultPreloads(db *gorm.DB) *gorm.DB {
	return db.Preload("Manager")
}

func Service(gctx golly.Context) EmployeeService {
	if service, ok := gctx.Get(serviceCtxKey); ok {
		return service.(EmployeeService)
	}

	service := DefaultEmployeeService{}
	gctx.Set(serviceCtxKey, service)

	return service
}
