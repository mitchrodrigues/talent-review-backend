package workos

import (
	"context"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
)

type CreateUserInput struct {
	FirstName     string
	LastName      string
	Password      string
	Email         string
	OrganzationID string
}

func (DefaultClient) CreateUser(ctx golly.Context, input CreateUserInput) (string, error) {
	response, err := usermanagement.CreateUser(ctx.Context(), usermanagement.CreateUserOpts{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  input.Password,
		Email:     input.Email,
	})

	if err != nil {
		return "", errors.WrapGeneric(err)
	}

	_, err = usermanagement.CreateOrganizationMembership(ctx.Context(), usermanagement.CreateOrganizationMembershipOpts{
		UserID:         response.ID,
		OrganizationID: input.OrganzationID,
	})

	return response.ID, errors.WrapGeneric(err)
}

func (DefaultClient) InviteUser(ctx golly.Context, organizationID, email, inviterID string) (string, error) {
	response, err := usermanagement.SendInvitation(
		context.Background(),
		usermanagement.SendInvitationOpts{
			Email:          email,
			InviterUserID:  inviterID,
			ExpiresInDays:  5,
			OrganizationID: organizationID,
		},
	)

	return response.ID, err
}
