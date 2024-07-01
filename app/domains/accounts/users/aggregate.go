package users

import (
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	es "github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
)

type Aggregate struct {
	eventsource.AggregateBase

	orm.ModelUUID

	FirstName string
	LastName  string
	Email     string

	OrganizationID uuid.UUID

	IdpID       string
	IdpInviteID string

	InvitedAt *time.Time
	InviterID uuid.UUID
}

func (*Aggregate) Topic() string                             { return "events.users" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return es.PostgresRepository{} }
func (*Aggregate) TableName() string                         { return "users" }

func (user *Aggregate) GetID() string   { return user.ID.String() }
func (user *Aggregate) SetID(id string) { user.ID, _ = uuid.Parse(id) }

func (user *Aggregate) RecordID() uuid.UUID { return user.ID }
func (user *Aggregate) RecordIdpID() string { return user.IdpID }

func (user *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	switch event := evt.Data.(type) {
	case UserCreated:
		user.ID = event.ID
		user.FirstName = event.FirstName
		user.OrganizationID = event.OrganizationID
		user.LastName = event.LastName
		user.Email = event.Email
		user.IdpID = event.IdpID

	case UserInvited:
		user.IdpInviteID = event.IdpInviteID
		user.InvitedAt = event.InvitedAt
		user.InviterID = event.InviterID

	case UserUpdated:
		user.IdpID = event.IdpID
		user.FirstName = event.FirstName
		user.LastName = event.LastName
		user.Email = event.Email
	}
}

var _ eventsource.Aggregate = &Aggregate{}
