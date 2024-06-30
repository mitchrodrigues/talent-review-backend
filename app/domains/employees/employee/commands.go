package employee

import (
	"fmt"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
)

type Create struct {
	Name  string
	Email string
	Title string

	OrganizationID uuid.UUID
	Manager        bool

	WorkerType EmployeeWorkerType
	Level      int

	TeamID uuid.UUID
}

func (cmd Create) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {

	fmt.Printf("\n\n%#v\n\n", cmd)

	id, _ := uuid.NewV7()

	tpe := IC
	if cmd.Manager {
		tpe = Manager
	}

	var userID *uuid.UUID
	if user, err := accounts.FindUserByEmail(ctx, cmd.Email); err != nil && user.ID != uuid.Nil {
		userID = &user.ID
	}

	eventsource.Apply(ctx, aggregate, Created{
		ID:             id,
		Name:           cmd.Name,
		Email:          cmd.Email,
		OrganizationID: cmd.OrganizationID,
		Level:          cmd.Level,
		Type:           tpe,
		WorkerType:     cmd.WorkerType,
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
