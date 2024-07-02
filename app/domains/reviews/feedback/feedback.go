package feedback

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

type Aggregate struct {
	eventsource.Aggregate

	orm.ModelUUID

	OwnerID        uuid.UUID
	EmployeeID     uuid.UUID
	OrganizationID uuid.UUID
	CycleID        uuid.UUID

	Email string
	Code  string

	Strength      string
	Opportunities string
	Additional    string
	Rating        int
}

func (*Aggregate) Topic() string                             { return "events.feedback" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return esbackend.PostgresRepository{} }
func (*Aggregate) TableName() string                         { return "feedbacks" }

func (feedback *Aggregate) GetID() string   { return feedback.ID.String() }
func (feedback *Aggregate) SetID(id string) { feedback.ID, _ = uuid.Parse(id) }

func (feedback *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	// switch event := evt.Data.(type) {
	// case TeamCreated:
	// }
}

var _ eventsource.Aggregate = &Aggregate{}
