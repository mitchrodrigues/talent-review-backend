package employees

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/passport"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
)

func EmployeeIDToIdentityMiddleware(next golly.HandlerFunc) golly.HandlerFunc {
	return func(c golly.WebContext) {
		ident := identity.FromContext(c.Context)

		id := Service(c.Context).PluckIDByUserID(c.Context, ident.UID)

		ident.EmployeeID = id

		c.Context = passport.ToContext(c.Context, ident)
		next(c)
	}
}
