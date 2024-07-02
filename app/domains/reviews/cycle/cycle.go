package cycle

import (
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

type Aggregate struct {
	eventsource.AggregateBase

	orm.ModelUUID

	OwnerID        uuid.UUID
	OrganizationID uuid.UUID

	StartAt time.Time
	EndAt   time.Time
}

func (*Aggregate) Topic() string                             { return "" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return esbackend.PostgresRepository{} }
func (*Aggregate) TableName() string                         { return "cycles" }

func (cycle *Aggregate) GetID() string   { return cycle.ID.String() }
func (cycle *Aggregate) SetID(id string) { cycle.ID, _ = uuid.Parse(id) }

func (cycle *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	switch event := evt.Data.(type) {
	case CycleCreated:
		cycle.ID = event.ID
		cycle.OrganizationID = event.OrganizationID
		cycle.OwnerID = event.OwnerID
		cycle.StartAt = event.StartAt
		cycle.EndAt = event.EndAt

		cycle.CreatedAt = evt.CreatedAt
		cycle.UpdatedAt = evt.CreatedAt

	}
}

var _ eventsource.Aggregate = &Aggregate{}
