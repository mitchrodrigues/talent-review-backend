package teams

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
)

type Create struct {
	Name           string
	LeadID         *uuid.UUID
	OrganizationID uuid.UUID
}

func (cmd Create) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	id, _ := uuid.NewV7()

	if cmd.OrganizationID == uuid.Nil {
		cmd.OrganizationID = identity.FromContext(gctx).OrganizationID
	}

	if cmd.LeadID != nil && *cmd.LeadID == uuid.Nil {
		cmd.LeadID = nil
	}

	eventsource.Apply(gctx, aggregate, Created{
		ID:             id,
		Name:           cmd.Name,
		LeadID:         cmd.LeadID,
		OrganizationID: cmd.OrganizationID,
	})

	return nil
}

type Update struct {
	Name   string
	LeadID *uuid.UUID
}

func (cmd Update) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	team := aggregate.(*Aggregate)

	eventsource.Apply(gctx, aggregate, Updated{
		Name:   helpers.Coalesce(cmd.Name, team.Name),
		LeadID: cmd.LeadID,
	})
	return nil
}
