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
	// Employees
	FindEmployeeByUserID(gctx golly.Context, userID uuid.UUID) (Employee, error)
	FindEmployeesByIDS(gctx golly.Context, ids uuid.UUIDs) ([]Employee, error)
	FindEmployeesForTeam(gctx golly.Context, teamID uuid.UUID, excludeEmployees ...uuid.UUID) ([]Employee, error)
	FindEmployeesByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) ([]Employee, error)
	FindEmployeesByManagerUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) ([]Employee, error)
	FindEmployeeByEmailAndOrganizationID(gctx golly.Context, email string, organizationID uuid.UUID) (Employee, error)
	FindEmployeesByManagerAndIDS(gctx golly.Context, managerID uuid.UUID, employeeIDs ...uuid.UUID) ([]Employee, error)
	FindEmployeeByID(gctx golly.Context, id uuid.UUID) (Employee, error)
	FindEmployeeEmailsBySearch(gctx golly.Context, name string) ([]string, error)

	PluckEmployeeIDsByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) (uuid.UUIDs, error)
	PluckIDByUserID(gctx golly.Context, userID uuid.UUID) uuid.UUID

	// Teams
	FindTeamsByOrganizationID(gctx golly.Context, organizationID uuid.UUID) ([]Team, error)
	FindTeamByID(gctx golly.Context, id uuid.UUID) (Team, error)

	// Roles
	FindRoleByID(gctx golly.Context, id uuid.UUID) (EmployeeRole, error)

	// Unsafe
	FindEmployeeByID_Unsafe(gctx golly.Context, id uuid.UUID) (Employee, error)
}

type DefaultEmployeeService struct{}

func baseEmployeeQuery(gctx golly.Context) *gorm.DB {
	return orm.DB(gctx).
		Model(&Employee{}).
		Where("employees.organization_id = ?", identity.FromContext(gctx).OrganizationID)
}

func (s DefaultEmployeeService) FindEmployeeByUserID(
	gctx golly.Context,
	userID uuid.UUID,
) (Employee, error) {
	var employee Employee

	err := baseEmployeeQuery(gctx).
		Where("user_id = ?", userID).
		First(&employee).
		Error

	return employee, errors.WrapNotFound(err)
}

func (s DefaultEmployeeService) FindEmployeesByIDS(
	gctx golly.Context,
	ids uuid.UUIDs,
) ([]Employee, error) {
	var employees []Employee

	err := baseEmployeeQuery(gctx).
		Where("id IN ?", ids).
		Find(&employees).
		Error

	return employees, err
}

func (s DefaultEmployeeService) FindEmployeesForTeam(
	gctx golly.Context,
	teamID uuid.UUID,
	excludeEmployees ...uuid.UUID,
) ([]Employee, error) {
	var employees []Employee

	query := baseEmployeeQuery(gctx).
		Where("team_id = ?", teamID)

	if len(excludeEmployees) > 0 {
		query = query.Where("id NOT IN ?", excludeEmployees)
	}

	err := query.Find(&employees).Error

	return employees, err
}

func (s DefaultEmployeeService) FindEmployeesByManagerID(
	gctx golly.Context,
	managerID uuid.UUID,
	scopes ...func(db *gorm.DB) *gorm.DB,
) ([]Employee, error) {
	var employees []Employee

	err := baseEmployeeQuery(gctx).
		Scopes(scopes...).
		Where("employees.manager_id = ?", managerID).
		Find(&employees).
		Error

	return employees, err
}

func (s DefaultEmployeeService) PluckIDByUserID(gctx golly.Context, userID uuid.UUID) uuid.UUID {
	var ids uuid.UUIDs

	baseEmployeeQuery(gctx).Where("user_id = ?", userID).
		Limit(1).
		Pluck("employees.id", &ids)

	if len(ids) == 0 {
		return uuid.Nil
	}

	return ids[0]
}

func (s DefaultEmployeeService) PluckEmployeeIDsByManagerID(
	gctx golly.Context,
	managerID uuid.UUID,
	scopes ...func(db *gorm.DB) *gorm.DB,
) (uuid.UUIDs, error) {
	var employeeIDs uuid.UUIDs

	err := baseEmployeeQuery(gctx).
		Scopes(scopes...).
		Where("employees.manager_id = ?", managerID).
		Pluck("employees.id", &employeeIDs).
		Error

	return employeeIDs, err
}

func (s DefaultEmployeeService) FindEmployeesByManagerUserID(
	gctx golly.Context,
	userID uuid.UUID,
	scopes ...func(db *gorm.DB) *gorm.DB,
) ([]Employee, error) {
	manager, err := s.FindEmployeeByUserID(gctx, userID)
	if err != nil || manager.ID == uuid.Nil {
		return nil, err
	}

	return s.FindEmployeesByManagerID(gctx, manager.ID, scopes...)
}

func (s DefaultEmployeeService) FindEmployeeByEmailAndOrganizationID(
	gctx golly.Context,
	email string,
	organizationID uuid.UUID,
) (Employee, error) {
	var employee Employee

	err := baseEmployeeQuery(gctx).
		Scopes(common.OrganizationIDScope(organizationID)).
		Where("email = ?", email).
		First(&employee).
		Error

	return employee, errors.WrapGeneric(err)
}

func (s DefaultEmployeeService) FindEmployeesByManagerAndIDS(
	gctx golly.Context,
	managerID uuid.UUID,
	employeeIDs ...uuid.UUID,
) ([]Employee, error) {
	var employees []Employee

	err := baseEmployeeQuery(gctx).
		Where("employees.manager_id = ?", managerID).
		Where("employees.id IN (?)", employeeIDs).
		Find(&employees).
		Error

	return employees, err
}

func (s DefaultEmployeeService) FindEmployeeByID(
	gctx golly.Context,
	id uuid.UUID,
) (Employee, error) {
	var employee Employee

	err := baseEmployeeQuery(gctx).
		Where("id = ?", id).
		First(&employee).
		Error

	return employee, errors.WrapNotFound(err)
}

func (s DefaultEmployeeService) FindEmployeeEmailsBySearch(
	gctx golly.Context,
	name string,
) ([]string, error) {
	var emails []string

	err := baseEmployeeQuery(gctx).
		Select("DISTINCT(email) AS email").
		Where("LOWER(email) LIKE ?", name+"%").
		Pluck("email", &emails).
		Error

	return emails, err
}

func (s DefaultEmployeeService) FindEmployeeByID_Unsafe(
	gctx golly.Context,
	id uuid.UUID,
) (Employee, error) {
	var employee Employee

	err := orm.DB(gctx).
		Model(&Employee{}).
		Where("id = ?", id).
		First(&employee).
		Error

	return employee, errors.WrapNotFound(err)
}

func (s DefaultEmployeeService) FindTeamsByOrganizationID(
	gctx golly.Context,
	organizationID uuid.UUID,
) ([]Team, error) {
	var teams []Team

	err := orm.DB(gctx).
		Model(&Team{}).
		Scopes(common.OrganizationIDScope(organizationID)).
		Scopes(s.DefaultPreloads).
		Find(&teams).
		Error

	return teams, err
}

func (s DefaultEmployeeService) FindTeamByID(
	gctx golly.Context,
	id uuid.UUID,
) (Team, error) {
	var team Team

	err := orm.DB(gctx).
		Model(&Team{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		First(&team, "id = ?", id).
		Error

	return team, err
}

func (s DefaultEmployeeService) FindRoleByID(gctx golly.Context, id uuid.UUID) (EmployeeRole, error) {
	var empRole EmployeeRole

	if id == uuid.Nil {
		return empRole, nil
	}

	err := orm.
		DB(gctx).
		Model(&empRole).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Find(&empRole, "id = ?", id).
		Error

	return empRole, err
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
