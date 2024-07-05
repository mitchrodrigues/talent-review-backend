package feedback

import (
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
)

type Create struct {
	OrganizationID uuid.UUID

	EmployeeID      uuid.UUID
	CollectionEndAt time.Time

	Email string
}

func (cmd Create) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	code, err := helpers.GenerateCode()
	if err != nil {
		return err
	}

	id, _ := uuid.NewV7()

	eventsource.Apply(gctx, aggregate, Created{
		ID:              id,
		Code:            code,
		Email:           cmd.Email,
		CollectionEndAt: cmd.CollectionEndAt,
		EmployeeID:      cmd.EmployeeID,
		OrganizationID:  cmd.OrganizationID,
	})

	return nil
}
