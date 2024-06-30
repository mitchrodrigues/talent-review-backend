package common

import (
	"github.com/golly-go/golly"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"gorm.io/gorm"
)

func OrganizationIDScopeForContext(gctx golly.Context) func(*gorm.DB) *gorm.DB {
	ident := identity.FromContext(gctx)

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("organization_id = ?", ident.OrganizationID)
	}
}

func OrganizationIDScope(orgID uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("organization_id = ?", orgID)
	}
}

func PaginationScope(page, limit int) func(*gorm.DB) *gorm.DB {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(offset)
	}
}
