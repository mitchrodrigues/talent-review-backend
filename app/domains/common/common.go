package common

import (
	"fmt"
	"strings"

	"github.com/golly-go/golly"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"gorm.io/gorm"
)

func OrganizationIDScopeForContext(gctx golly.Context, tablePrefix ...string) func(*gorm.DB) *gorm.DB {
	ident := identity.FromContext(gctx)

	table := ""
	if len(tablePrefix) > 0 {
		table = tablePrefix[0] + "."
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%sorganization_id = ?", table), ident.OrganizationID)
	}
}

func OrganizationIDScope(orgID uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("organization_id = ?", orgID)
	}
}

func EmailScope(email string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(email) = ?", strings.ToLower(email))
	}
}
