package cycle

import (
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
)

type FindOrCreateCycle struct {
	Type    string
	StartAt time.Time
	EndAt   time.Time
}

func (cmd FindOrCreateCycle) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {

	ident := identity.FromContext(gctx)

	// Find the cycle that matches the criteria
	err := orm.DB(gctx).
		Model(&Aggregate{}).
		Scopes(common.OrganizationIDScope(ident.OrganizationID)).
		Where("end_at > ?", cmd.StartAt).                               // Only consider cycles that haven't ended before the new cycle starts
		Where("(start_at < ? AND end_at > ?)", cmd.EndAt, cmd.StartAt). // Check for overlap
		Where("organization_id = ?", ident.OrganizationID).
		Where("owner_id = ?", ident.UID).
		First(aggregate).
		Error

	if err == nil {
		return nil
	}

	id, _ := uuid.NewV7()

	eventsource.Apply(gctx, aggregate, CycleCreated{
		ID:             id,
		StartAt:        cmd.StartAt,
		EndAt:          cmd.EndAt,
		OwnerID:        ident.UID,
		OrganizationID: ident.OrganizationID,
		Type:           cmd.Type,
	})
	return nil
}
