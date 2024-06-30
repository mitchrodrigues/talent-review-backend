package teams

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
)

type CreateTeam struct {
	Name           string
	ManagerID      uuid.UUID
	OrganizationID uuid.UUID
}

func (cmd CreateTeam) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	id, _ := uuid.NewV7()

	if cmd.OrganizationID == uuid.Nil {
		cmd.OrganizationID = identity.FromContext(gctx).OrganizationID
	}

	eventsource.Apply(gctx, aggregate, TeamCreated{
		ID:             id,
		Name:           cmd.Name,
		ManagerID:      cmd.ManagerID,
		OrganizationID: cmd.OrganizationID,
	})

	return nil
}
