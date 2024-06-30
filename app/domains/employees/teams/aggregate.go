package teams

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

type Aggregate struct {
	eventsource.AggregateBase

	orm.ModelUUID

	Name string

	ManagerID      uuid.UUID
	OrganizationID uuid.UUID
}

func (*Aggregate) Topic() string                             { return "events.teams" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return esbackend.PostgresRepository{} }
func (*Aggregate) TableName() string                         { return "teams" }

func (team *Aggregate) GetID() string   { return team.ID.String() }
func (team *Aggregate) SetID(id string) { team.ID, _ = uuid.Parse(id) }

func (team *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	switch event := evt.Data.(type) {
	case TeamCreated:
		team.ID = event.ID
		team.ManagerID = event.ManagerID
		team.OrganizationID = event.OrganizationID
		team.Name = event.Name

		team.CreatedAt = evt.CreatedAt
		team.UpdatedAt = evt.CreatedAt
	}
}

var _ eventsource.Aggregate = &Aggregate{}
