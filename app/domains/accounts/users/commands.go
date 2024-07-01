package users

import (
	"fmt"
	"strings"
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

	Name  string
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

	idpID, url, err := cmd.WorkosClient.InviteUser(ctx, cmd.Organization.RecordIdpID(), cmd.Email, inviterIdpID)
	if err != nil {
		return err
	}

	var firstName string
	var lastName string

	pieces := strings.Split(cmd.Name, " ")
	switch len(pieces) {
	case 0:
	case 1:
		firstName = pieces[0]
	case 2:
		firstName = pieces[0]
		lastName = pieces[1]
	default:
		firstName = strings.Join(pieces[0:len(pieces)-2], " ")
		lastName = pieces[len(pieces)-1]
	}

	t := time.Now()

	eventsource.Apply(ctx, aggregate, UserCreated{
		ID:             id,
		FirstName:      firstName,
		LastName:       lastName,
		Email:          cmd.Email,
		OrganizationID: cmd.Organization.RecordID(),
	})

	eventsource.Apply(ctx, aggregate, UserInvited{
		IdpInviteID: idpID,
		InvitedAt:   &t,
		InviterID:   inviterID,
		InviteURL:   url,
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
		Email:     helpers.Coalesce(cmd.Email, user.Email),
	})

	return nil
}
