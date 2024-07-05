package employee

import (
	"fmt"
	"strings"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
)

type Create struct {
	Name  string `validate:"required"`
	Email string `validate:"email"`
	Title string

	OrganizationID uuid.UUID `validate:"required"`
	Manager        bool

	WorkerType EmployeeWorkerType
	Level      int `validate:"lt=10"`

	TeamID uuid.UUID
}

func (cmd Create) Validate(ctx golly.Context, aggregate eventsource.Aggregate) error {
	var employee Aggregate

	orm.DB(ctx).
		Scopes(common.EmailScope(cmd.Email)).
		Find(&employee)

	if employee.ID != uuid.Nil {
		return fmt.Errorf("employee with that email already exists")
	}

	return nil
}

func (cmd Create) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	id, _ := uuid.NewV7()

	tpe := IC
	if cmd.Manager {
		tpe = Manager
	}

	var userID *uuid.UUID
	if user, err := accounts.FindUserByEmail(ctx, cmd.Email); err != nil && user.ID != uuid.Nil {
		userID = &user.ID
	}

	workerType := EmployeeWorkerType(strings.TrimSpace(string(cmd.WorkerType)))

	eventsource.Apply(ctx, aggregate, Created{
		ID:             id,
		Name:           cmd.Name,
		Email:          cmd.Email,
		OrganizationID: cmd.OrganizationID,
		Level:          cmd.Level,
		Type:           tpe,
		WorkerType:     workerType,
		UserID:         userID,
	})

	if cmd.TeamID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, TeamUpdated{cmd.TeamID})
	}

	if cmd.Title != "" {
		eventsource.Apply(ctx, aggregate, TitleUpdated{cmd.Title})
	}

	return nil
}

type Update struct {
	Email string
	Name  string
	Level int

	Title      string
	WorkerType EmployeeWorkerType
	TeamID     uuid.UUID
}

func (cmd Update) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	employee := aggregate.(*Aggregate)

	level := employee.Level
	if cmd.Level != 0 {
		level = cmd.Level
	}

	workerType := EmployeeWorkerType(helpers.Coalesce(strings.TrimSpace(string(cmd.WorkerType)), strings.TrimSpace(string(employee.WorkerType))))

	eventsource.Apply(ctx, aggregate, Updated{
		Name:       helpers.Coalesce(cmd.Name, employee.Name),
		Email:      helpers.Coalesce(cmd.Email, employee.Email),
		Level:      level,
		WorkerType: workerType,
	})

	if cmd.TeamID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, TeamUpdated{cmd.TeamID})
	}

	if cmd.Title != "" {
		eventsource.Apply(ctx, aggregate, TitleUpdated{cmd.Title})
	}

	return nil
}

type UpdateUser struct {
	UserID uuid.UUID
}

func (cmd UpdateUser) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	eventsource.Apply(gctx, aggregate, UserUpdated(cmd))

	return nil
}

type UpdateTeam struct {
	TeamID uuid.UUID
}

func (cmd UpdateTeam) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	eventsource.Apply(gctx, aggregate, TeamUpdated(cmd))

	return nil
}
