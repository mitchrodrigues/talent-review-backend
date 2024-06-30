package employee

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

// id UUID NOT NULL,
// email VARCHAR(255) NOT NULL,
// name VARCHAR(255) NOT NULL,
// title VARCHAR(255) NOT NULL,
// type char(4) NOT NULL,
// level char(4) NOT NULL,
// team_id UUID,
// level_start_at TIMESTAMP,
// employement_start_at TIMESTAMP,
// organization_id UUID NOT NULL,
// created_at TIMESTAMP,
// updated_at TIMESTAMP,
// deleted_at TIMESTAMP,
// PRIMARY KEY(id)

type EmployeeType string

type EmployeeLevel string

const (
	IC      EmployeeType = "IC"
	Manager EmployeeType = "MNG"
)

type EmployeeWorkerType string

const (
	DirectContractor EmployeeWorkerType = "DC"
	AgencyContractor EmployeeWorkerType = "AC"
	FTE              EmployeeWorkerType = "FTE"
)

type Aggregate struct {
	eventsource.AggregateBase

	orm.ModelUUID

	Name  string
	Email string
	Title string

	OrganizationID uuid.UUID
	UserID         *uuid.UUID
	TeamID         *uuid.UUID

	Level      int
	Type       EmployeeType
	WorkerType EmployeeWorkerType
}

func (*Aggregate) Topic() string                             { return "events.employees" }
func (*Aggregate) TableName() string                         { return "employees" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return esbackend.PostgresRepository{} }

func (employee *Aggregate) GetID() string   { return employee.ID.String() }
func (employee *Aggregate) SetID(id string) { employee.ID, _ = uuid.Parse(id) }

func (employee *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	switch event := evt.Data.(type) {
	case Created:
		employee.ID = event.ID
		employee.Name = event.Name
		employee.Email = event.Email
		employee.Level = event.Level
		employee.Type = event.Type
		employee.WorkerType = event.WorkerType

		employee.OrganizationID = event.OrganizationID
		employee.UserID = event.UserID

		employee.CreatedAt = evt.CreatedAt
		employee.UpdatedAt = evt.CreatedAt

	case TitleUpdated:
		employee.Title = event.Title
		employee.UpdatedAt = evt.CreatedAt

	case UserUpdated:
		employee.UserID = &event.UserID
		employee.UpdatedAt = evt.CreatedAt

	case TeamUpdated:
		employee.TeamID = &event.TeamID
		employee.UpdatedAt = evt.CreatedAt
	}
}

var _ eventsource.Aggregate = &Aggregate{}
