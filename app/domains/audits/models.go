package audits

import (
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

type Event struct {
	esbackend.Event

	User accounts.User
}
