package role

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

type EmployeeType string

const (
	IC      EmployeeType = "IC"
	Manager EmployeeType = "MNG"
)

type Aggregate struct {
	eventsource.AggregateBase

	orm.ModelUUID

	Title          string
	Level          int
	OrganizationID uuid.UUID
	Track          EmployeeType
}

func (*Aggregate) Topic() string                             { return "events.employee_roles" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return esbackend.PostgresRepository{} }
func (*Aggregate) TableName() string                         { return "employee_roles" }

func (role *Aggregate) GetID() string   { return role.ID.String() }
func (role *Aggregate) SetID(id string) { role.ID, _ = uuid.Parse(id) }

func (role *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	switch event := evt.Data.(type) {
	case Created:
		role.ID = event.ID
		role.OrganizationID = event.OrganizationID
		role.Title = event.Title
		role.Level = event.Level
		role.Track = event.Track

		role.CreatedAt = evt.CreatedAt

	case TitleUpdated:
		role.Title = event.Title

	case LevelUpdated:
		role.Level = event.Level

	case TrackUpdated:
		role.Track = event.Track

	}
	role.UpdatedAt = evt.CreatedAt
}

var _ eventsource.Aggregate = &Aggregate{}
