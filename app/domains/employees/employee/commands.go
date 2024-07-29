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

	OrganizationID uuid.UUID `validate:"required"`

	WorkerType     EmployeeWorkerType `validate:"required"`
	TeamID         uuid.UUID
	ManagerID      uuid.UUID
	EmployeeRoleID uuid.UUID
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

	eventsource.Apply(ctx, aggregate, Created{
		ID:             id,
		Name:           cmd.Name,
		Email:          cmd.Email,
		OrganizationID: cmd.OrganizationID,
	})

	if user, err := accounts.FindUserByEmail(ctx, cmd.Email); err != nil && user.ID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, UserUpdated{user.ID})
	}

	if cmd.WorkerType != "" {
		eventsource.Apply(ctx, aggregate, WorkerTypeUpdated{EmployeeWorkerType(strings.TrimSpace(string(cmd.WorkerType)))})
	}

	if cmd.EmployeeRoleID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, RoleUpdated{cmd.EmployeeRoleID})
	}

	if cmd.ManagerID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, ManagerUpdated{&cmd.ManagerID})
	}

	if cmd.TeamID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, TeamUpdated{&cmd.TeamID})
	}

	return nil
}

type Update struct {
	Email string
	Name  string

	WorkerType EmployeeWorkerType
	TeamID     uuid.UUID
	ManagerID  uuid.UUID

	EmployeeRoleID uuid.UUID
}

func (cmd Update) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	employee := aggregate.(*Aggregate)

	if cmd.Name != "" || cmd.Email != "" {
		eventsource.Apply(ctx, aggregate, PersonalDetailsUpdated{
			Name:  helpers.Coalesce(cmd.Name, employee.Name),
			Email: helpers.Coalesce(cmd.Email, employee.Email),
		})
	}

	if cmd.WorkerType != "" {
		eventsource.Apply(ctx, aggregate, EmployeeWorkerType(strings.TrimSpace(string(cmd.WorkerType))))
	}

	if cmd.TeamID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, TeamUpdated{&cmd.TeamID})
	}

	if cmd.ManagerID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, ManagerUpdated{&cmd.ManagerID})
	}

	if cmd.EmployeeRoleID != uuid.Nil {
		eventsource.Apply(ctx, aggregate, RoleUpdated{cmd.EmployeeRoleID})
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
