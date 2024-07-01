package mailgun

import (
	"fmt"

	"github.com/golly-go/golly"
)

type InviteEmailParams struct {
	Name             string
	InvitorName      string
	Email            string
	AcceptLink       string
	OrganizationName string
}

func (c *DefaultClient) SendInviteEmail(gctx golly.Context, params InviteEmailParams) error {
	return c.SendEmailTemplate(gctx, EmailWithTemplate{
		Email: Email{
			Recipient: params.Email,
			Subject:   fmt.Sprintf("You have been invited to join %s on TalentRadar", params.OrganizationName),
		},
		Template: "user-invite",
		Variables: map[string]interface{}{
			"name":         params.Name,
			"invitorName":  params.InvitorName,
			"email":        params.Email,
			"acceptLink":   params.AcceptLink,
			"organization": params.OrganizationName,
		},
	})
}
