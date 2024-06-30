package accounts

import (
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts/organizations"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts/users"
)

type User struct {
	users.Aggregate

	Organization Organization
}

type Organization struct {
	organizations.Aggregate

	Users []User
}
