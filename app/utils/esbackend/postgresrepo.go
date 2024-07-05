package esbackend

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"gorm.io/gorm"
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
		event, err := mapToDB(ctx, t)
		if err != nil {
			return err
		}
		return orm.NewDB(ctx).Model(event).Create(&event).Error
	default:
		return orm.NewDB(ctx).Model(t).Session(&gorm.Session{FullSaveAssociations: true}).Save(t).Error
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

func mapToDB(gctx golly.Context, evt *eventsource.Event) (Event, error) {
	var err error

	agID, _ := uuid.Parse(evt.AggregateID)

	ident := identity.FromContext(gctx)

	// organizationID := ident.OrganizationID
	// if organizationID == uuid.Nil {
	// 	if oid, err := GetOrganizationID(evt.Data); err != nil {
	// 		organizationID = oid
	// 	}
	// }

	ret := Event{
		ModelUUID: orm.ModelUUID{
			ID:        evt.ID,
			CreatedAt: evt.CreatedAt,
		},
		Type:           evt.Event,
		AggregateID:    agID,
		AggregateType:  evt.AggregateType,
		UserID:         &ident.UID,
		OrganizationID: &ident.OrganizationID,
		Version:        evt.Version,
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

// GetOrganizationID extracts the OrganizationID field from any struct using reflection.
func GetOrganizationID(obj interface{}) (uuid.UUID, error) {
	v := reflect.ValueOf(obj)

	// Check if the input is a pointer and get the underlying element.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ensure we have a struct.
	if v.Kind() != reflect.Struct {
		return uuid.Nil, fmt.Errorf("expected a struct, but got %v", v.Kind())
	}

	// Get the OrganizationID field.
	field := v.FieldByName("OrganizationID")
	if !field.IsValid() {
		return uuid.Nil, fmt.Errorf("field OrganizationID not found in struct")
	}

	// Ensure the field is of the correct type.
	if field.Type() != reflect.TypeOf(uuid.UUID{}) {
		return uuid.Nil, fmt.Errorf("expected OrganizationID to be uuid.UUID, but got %v", field.Type())
	}

	return field.Interface().(uuid.UUID), nil
}
