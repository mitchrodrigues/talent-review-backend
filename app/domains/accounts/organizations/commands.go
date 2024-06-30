package organizations

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
)

type CreateOrganization struct {
	workos.WorkosClient

	Name string
}

func (cmd CreateOrganization) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	id, _ := uuid.NewV7()

	idpID, err := cmd.WorkosClient.CreateOrganization(ctx, cmd.Name)
	if err != nil {
		return err
	}

	eventsource.Apply(ctx, aggregate, OrganizationCreated{
		ID:    id,
		Name:  cmd.Name,
		IdpID: idpID,
	})

	return nil
}
