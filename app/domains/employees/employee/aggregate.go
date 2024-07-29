package employee

import (
	"time"

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

type EmployeeLevel string

type EmployeeWorkerType string

const (
	DirectContractor EmployeeWorkerType = "DC"
	AgencyContractor EmployeeWorkerType = "AC"
	FTE              EmployeeWorkerType = "FTE"
)

type EmployeeHistory struct {
	orm.ModelUUID

	EmployeeID uuid.UUID  `gorm:"type:uuid;not null"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null"`
	Change     ChangeData `gorm:"type:jsonb;not null"`
}

type ChangeData struct {
	Previous interface{} `json:"previous"`
	Current  interface{} `json:"current"`
	Field    string      `json:"field"`
}

type Aggregate struct {
	eventsource.AggregateBase

	orm.ModelUUID

	Name  string
	Email string

	OrganizationID uuid.UUID
	ManagerID      *uuid.UUID
	UserID         *uuid.UUID
	TeamID         *uuid.UUID

	WorkerType EmployeeWorkerType

	EmployeeRoleID uuid.UUID

	TerminatedAt *time.Time
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
		employee.OrganizationID = event.OrganizationID

		employee.CreatedAt = evt.CreatedAt

	case UserUpdated:
		employee.UserID = &event.UserID

	case TeamUpdated:
		employee.TeamID = event.TeamID

	case ManagerUpdated:
		employee.ManagerID = event.ManagerID

	case Updated:
		employee.Name = event.Name
		employee.Email = event.Email

	case WorkerTypeUpdated:
		employee.WorkerType = event.WorkerType

	case RoleUpdated:
		employee.EmployeeRoleID = event.EmployeeRoleID

	case Terminate:
		employee.TerminatedAt = &event.TerminatedAt
	}

	employee.UpdatedAt = evt.CreatedAt
}

var _ eventsource.Aggregate = &Aggregate{}
