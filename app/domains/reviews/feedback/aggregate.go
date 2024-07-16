package feedback

import (
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

type FeedbackDetails struct {
	orm.ModelUUID

	EmployeeID     uuid.UUID
	FeedbackID     uuid.UUID
	OrganizationID uuid.UUID

	Strengths     string
	Opportunities string
	Additional    string
	EnoughData    bool

	Rating int
}

type FeedbackSummary struct {
	orm.ModelUUID

	EmployeeID     uuid.UUID
	FeedbackID     uuid.UUID
	OrganizationID uuid.UUID

	Summary     string
	ActionItems string
}

type Aggregate struct {
	eventsource.AggregateBase

	orm.ModelUUID

	OwnerID        uuid.UUID
	EmployeeID     uuid.UUID
	OrganizationID uuid.UUID

	Email string
	Code  string

	SubmittedAt     *time.Time
	CollectionEndAt time.Time

	Details FeedbackDetails `gorm:"foreignKey:FeedbackID"`
	Summary FeedbackSummary `gorm:"foreignKey:FeedbackID"`
}

func (*Aggregate) Topic() string                             { return "events.feedback" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return esbackend.PostgresRepository{} }
func (*Aggregate) TableName() string                         { return "feedbacks" }

func (feedback *Aggregate) GetID() string   { return feedback.ID.String() }
func (feedback *Aggregate) SetID(id string) { feedback.ID, _ = uuid.Parse(id) }

func (feedback *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	switch event := evt.Data.(type) {
	case Created:
		feedback.ID = event.ID
		feedback.Email = event.Email
		feedback.CollectionEndAt = event.CollectionEndAt
		feedback.EmployeeID = event.EmployeeID
		feedback.Code = event.Code
		feedback.OrganizationID = event.OrganizationID

		feedback.OwnerID = event.OwnerID
		feedback.CreatedAt = evt.CreatedAt
		feedback.UpdatedAt = evt.CreatedAt

	case DetailsCreated:
		feedback.Details.ID = event.ID
		feedback.Details.FeedbackID = event.FeedbackID
		feedback.Details.OrganizationID = event.OrganizationID
		feedback.Details.EmployeeID = event.EmployeeID

	case DetailsUpdated:
		feedback.Details.Strengths = event.Strenghts
		feedback.Details.Opportunities = event.Opportunities
		feedback.Details.Rating = event.Rating
		feedback.Details.Additional = event.Additional
		feedback.Details.EnoughData = event.EnoughData
		feedback.Details.CreatedAt = evt.CreatedAt
		feedback.Details.UpdatedAt = evt.CreatedAt

	case SummaryCreated:
		feedback.Summary.ID = event.ID
		feedback.Summary.FeedbackID = event.FeedbackID
		feedback.Summary.OrganizationID = event.OrganizationID
		feedback.Summary.EmployeeID = event.EmployeeID

	case SummaryUpdated:
		feedback.Summary.Summary = event.Summary
		feedback.Summary.ActionItems = event.ActionItems

	case Submitted:
		feedback.SubmittedAt = &evt.CreatedAt

	}
}

var _ eventsource.Aggregate = &Aggregate{}
