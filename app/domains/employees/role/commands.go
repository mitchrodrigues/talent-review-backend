package role

import (
	"fmt"
	"strings"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
)

var (
	ErrorRoleExists = fmt.Errorf("role exists")
)

type Create struct {
	OrganizationID uuid.UUID    `validate:"uuid"`
	Title          string       `validate:"required"`
	Level          int          `validate:"gte=0,lte=11"`
	Track          EmployeeType `validate:"required"`
}

func (cmd Create) Validate(ctx golly.Context, aggregate eventsource.Aggregate) error {
	var role Aggregate

	orm.DB(ctx).Model(&role).Find(&role, map[string]interface{}{
		"title":           cmd.Title,
		"organization_id": cmd.OrganizationID,
	})

	if role.ID != uuid.Nil {
		return errors.WrapInvalidFields(ErrorRoleExists)
	}

	return nil

}

func (cmd Create) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	id, _ := uuid.NewV7()

	eventsource.Apply(ctx, aggregate, Created{
		ID:             id,
		OrganizationID: cmd.OrganizationID,
		Title:          cmd.Title,
		Level:          cmd.Level,
		Track:          cmd.Track,
	})
	return nil
}

type Update struct {
	Title string
	Track string
	Level int
}

func (cmd Update) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {

	fmt.Printf("%#v\n", cmd)

	if cmd.Title != "" {
		eventsource.Apply(ctx, aggregate, TitleUpdated{cmd.Title})
	}

	if cmd.Level != 0 {
		eventsource.Apply(ctx, aggregate, LevelUpdated{cmd.Level})
	}

	if cmd.Track != "" {
		track := IC

		switch strings.ToLower(cmd.Track) {
		case "manager", "mng":
			track = Manager
		}

		eventsource.Apply(ctx, aggregate, TrackUpdated{EmployeeType(track)})
	}

	return nil
}
