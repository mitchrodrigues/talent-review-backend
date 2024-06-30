package accounts

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindUserByEmail(gctx golly.Context, email string) (User, error) {
	var user User

	err := orm.DB(gctx).Model(user).Find(&user, "email = ?", email).Error

	return user, err
}

func FindUserByID(gctx golly.Context, id string, scopes ...func(*gorm.DB) *gorm.DB) (User, error) {
	var user User

	err := orm.DB(gctx).Model(user).Scopes(scopes...).Find(&user, "id = ?", id).Error
	if user.ID == uuid.Nil {
		return user, errors.WrapNotFound(gorm.ErrRecordNotFound)
	}

	return user, err
}

func FindUserByIDPId(gctx golly.Context, idpID string) (User, error) {
	var user User

	err := orm.DB(gctx).Model(user).Find(&user, "idp_id = ?", idpID).Error

	return user, err
}

func FindOrganizationByID(gctx golly.Context, id uuid.UUID) (Organization, error) {
	var organization Organization

	err := orm.
		NewDB(gctx).
		Model(&organization).
		Find(&organization, "id = ?", id).
		Error

	return organization, err
}

func FindOrganizationByIDPId(ctx golly.Context, idpID string) (Organization, error) {
	var organization Organization

	err := orm.DB(ctx).Model(organization).Find(&organization, "idp_id = ?", idpID).Error
	return organization, err
}

func DefaultUserPreloads(db *gorm.DB) *gorm.DB {
	return db.Preload("Organization")
}
