package organizations

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

	Name  string
	IdpID string

	MerchantCustomerID string
	MerchantPlanID     string
	MerchantPlanName   string

	ActivatedAt   *time.Time
	DeactivatedAt *time.Time
}

func (*Aggregate) Topic() string                             { return "events.organization" }
func (*Aggregate) Repo(golly.Context) eventsource.Repository { return es.PostgresRepository{} }
func (*Aggregate) TableName() string                         { return "organizations" }

func (org *Aggregate) RecordID() uuid.UUID { return org.ID }
func (org *Aggregate) RecordIdpID() string { return org.IdpID }

func (org *Aggregate) GetID() string   { return org.ID.String() }
func (org *Aggregate) SetID(id string) { org.ID, _ = uuid.Parse(id) }

func (org *Aggregate) Apply(ctx golly.Context, evt eventsource.Event) {
	switch event := evt.Data.(type) {
	case OrganizationCreated:
		org.ID = event.ID
		org.Name = event.Name
		org.IdpID = event.IdpID
		org.MerchantPlanName = event.PlanName
	}
}

var _ eventsource.Aggregate = &Aggregate{}
