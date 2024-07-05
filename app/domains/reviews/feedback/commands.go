package feedback

import (
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
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
		OwnerID:         identity.FromContext(gctx).UID,
	})

	return nil
}

type Submit struct{}

func (Submit) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	eventsource.Apply(gctx, aggregate, Submitted{})
	return nil
}

type CreateOrUpdateDetails struct {
	Strength      string
	Opportunities string
	Additional    string

	Rating int

	EnoughData *bool
}

func (cmd CreateOrUpdateDetails) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	feedback := aggregate.(*Aggregate)

	orm.
		DB(gctx).
		Model(feedback.Details).
		Find(&feedback.Details, "feedback_id = ?", feedback.ID)

	if feedback.Details.ID == uuid.Nil {
		id, _ := uuid.NewV7()
		eventsource.Apply(gctx, aggregate, DetailsCreated{
			ID:             id,
			OrganizationID: feedback.OrganizationID,
			EmployeeID:     feedback.EmployeeID,
			FeedbackID:     feedback.ID,
		})
	}

	rating := feedback.Details.Rating
	if cmd.Rating > 0 {
		rating = cmd.Rating
	}

	enoughData := feedback.Details.EnoughData
	if cmd.EnoughData != nil {
		enoughData = *cmd.EnoughData
	}

	eventsource.Apply(gctx, aggregate, DetailsUpdated{
		Strenghts:     helpers.Coalesce(cmd.Strength, feedback.Details.Strenghts),
		Opportunities: helpers.Coalesce(cmd.Opportunities, feedback.Details.Opportunities),
		Additional:    helpers.Coalesce(cmd.Additional, feedback.Details.Additional),
		Rating:        rating,
		EnoughData:    enoughData,
	})

	return nil
}
