package esbackend

import (
	"encoding/json"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresRepository struct{}

type Event struct {
	orm.ModelUUID

	AggregateID   uuid.UUID `json:"aggregateID"`
	AggregateType string    `json:"aggregateType"`

	Version uint   `json:"version"`
	Type    string `json:"type"`

	RawData     postgres.Jsonb `json:"-" gorm:"type:jsonb;column:data"`
	RawMetadata postgres.Jsonb `json:"-" gorm:"type:jsonb;column:metadata"`

	OrganizationID *uuid.UUID `json:"organizationID"`
	UserID         *uuid.UUID `json:"userID"`
}

func (pr PostgresRepository) Load(ctx golly.Context, object interface{}) error {
	return orm.DB(ctx).Model(object).First(object).Error
}

func (PostgresRepository) Save(ctx golly.Context, object interface{}) error {
	switch t := object.(type) {
	case *eventsource.Event:
		event, err := mapToDB(t)
		if err != nil {
			return err
		}
		return orm.NewDB(ctx).Model(event).Create(&event).Error
	default:
		return orm.NewDB(ctx).Model(t).Save(t).Error
	}
}

func (PostgresRepository) IsNewRecord(obj interface{}) bool {
	if ag, ok := obj.(eventsource.Aggregate); ok {
		return ag.GetID() == uuid.Nil.String()
	}
	return false
}

func (r PostgresRepository) Transaction(ctx golly.Context, fn func(golly.Context, eventsource.Repository) error) error {
	return fn(ctx, r)
}

func mapToDB(evt *eventsource.Event) (Event, error) {
	var err error

	agID, _ := uuid.Parse(evt.AggregateID)

	ret := Event{
		ModelUUID: orm.ModelUUID{
			ID:        evt.ID,
			CreatedAt: evt.CreatedAt,
		},
		Type:          evt.Event,
		AggregateID:   agID,
		AggregateType: evt.AggregateType,
		Version:       evt.Version,
	}

	ret.RawMetadata.RawMessage, err = json.Marshal(evt.Metadata)
	if err != nil {
		return ret, err
	}

	ret.RawData.RawMessage, err = json.Marshal(evt.Data)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
