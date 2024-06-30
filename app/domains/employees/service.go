package employees

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"gorm.io/gorm"
)

// func EmployeeByOrganition(gctx golly.Context, organizationID uuid.UUID, page, limit int)

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
		Find(&emp, "id = ?", id).
		Error

	return emp, errors.WrapNotFound(err)
}

func FindTeamsByOrganizationID(gctx golly.Context, organizationID uuid.UUID) ([]Team, error) {
	var teams []Team

	err := orm.
		DB(gctx).
		Model(Team{}).
		Scopes(DefaultPreloads).
		Find(&teams).
		Error

	return teams, err
}

func DefaultPreloads(db *gorm.DB) *gorm.DB {
	return db.Preload("Manager")
}
