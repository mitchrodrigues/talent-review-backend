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

func UserIsManagerScope(gctx golly.Context, table string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ident := identity.FromContext(gctx)
		return db.Joins(fmt.Sprintf("JOIN employees employee ON %s.employee_id = employee.id", table)).
			Joins("JOIN employees manager ON manager.user_id = ?", ident.UID).
			Joins("JOIN teams team ON employee.team_id IS NOT NULL AND team.id = employee.team_id AND team.manager_id = manager.id")
	}
}
