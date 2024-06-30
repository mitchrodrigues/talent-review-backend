package users

import (
	"fmt"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
)

type CreateUser struct {
	workos.WorkosClient

	Organization workos.IDPObject

	Email     string
	FirstName string
	LastName  string
	Password  string
}

func (cmd CreateUser) Validate(golly.Context, eventsource.Aggregate) error {
	if cmd.WorkosClient == nil {
		return errors.WrapFatal(fmt.Errorf("workosclient not provided"))
	}
	return nil
}

func (cmd CreateUser) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	id, _ := uuid.NewV7()

	idpID, err := cmd.WorkosClient.CreateUser(ctx, workos.CreateUserInput{
		FirstName:     cmd.FirstName,
		LastName:      cmd.LastName,
		Email:         cmd.Email,
		Password:      cmd.Password,
		OrganzationID: cmd.Organization.RecordIdpID(),
	})

	if err != nil {
		return err
	}

	eventsource.Apply(ctx, aggregate, UserCreated{
		ID:             id,
		IdpID:          idpID,
		FirstName:      cmd.FirstName,
		LastName:       cmd.LastName,
		Email:          cmd.Email,
		OrganizationID: cmd.Organization.RecordID(),
	})

	return nil
}

var _ eventsource.Command = CreateUser{}

type InviteUser struct {
	workos.WorkosClient

	Organization workos.IDPObject
	Inviter      workos.IDPObject

	Email string
}

func (cmd InviteUser) Validate(golly.Context, eventsource.Aggregate) error {
	if cmd.Organization == nil {
		return errors.WrapUnprocessable(fmt.Errorf("organization is required"))
	}

	if cmd.WorkosClient == nil {
		return errors.WrapFatal(fmt.Errorf("workosclient not provided"))
	}

	return nil
}

func (cmd InviteUser) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	id, _ := uuid.NewV7()
	inviterIdpID := ""
	inviterID := uuid.Nil

	if cmd.Inviter != nil {
		inviterIdpID = cmd.Inviter.RecordIdpID()
		inviterID = cmd.Inviter.RecordID()
	}

	idpID, err := cmd.WorkosClient.InviteUser(ctx, cmd.Organization.RecordIdpID(), cmd.Email, inviterIdpID)
	if err != nil {
		return err
	}

	t := time.Now()

	eventsource.Apply(ctx, aggregate, UserCreated{
		ID:             id,
		Email:          cmd.Email,
		IdpInviteID:    idpID,
		OrganizationID: cmd.Organization.RecordID(),

		InvitedAt: &t,
		InviterID: inviterID,
	})

	return nil
}

var _ eventsource.Command = InviteUser{}

type EditUser struct {
	IdpID string

	Email          string
	FirstName      string
	LastName       string
	ProfilePicture string
}

func (cmd EditUser) Perform(ctx golly.Context, aggregate eventsource.Aggregate) error {
	user := aggregate.(*Aggregate)

	eventsource.Apply(ctx, aggregate, UserUpdated{
		IdpID:     helpers.Coalesce(cmd.IdpID, user.IdpID),
		FirstName: helpers.Coalesce(cmd.FirstName, user.FirstName),
		LastName:  helpers.Coalesce(cmd.LastName, user.LastName),
		// ProfilePicture: helpers.Coalesce(cmd.ProfilePicture, user.ProfilePicture),
		Email: helpers.Coalesce(cmd.Email, user.Email),
	})

	return nil
}
