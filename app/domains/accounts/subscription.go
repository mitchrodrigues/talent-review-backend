package accounts

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts/users"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/mailgun"
)

func SendInviteEmail(gctx golly.Context, agg eventsource.Aggregate, evt eventsource.Event) error {

	switch event := evt.Data.(type) {
	case users.UserInvited:
		go func(gctx golly.Context, agg eventsource.Aggregate, evt eventsource.Event) {
			user := agg.(*users.Aggregate)

			invitorName := "a coworker"

			if event.InviterID != uuid.Nil {
				invitor, err := FindUserByID(gctx, event.InviterID.String())
				if err != nil {
					return
				}

				invitorName = invitor.FirstName
			}

			org, err := FindOrganizationByID(gctx, user.OrganizationID)
			if err != nil {
				gctx.Logger().Warnf("Unable to send invite email to %s %v", evt.AggregateID, err)
				return
			}

			err = mailgun.GetClient(gctx).SendInviteEmail(gctx, mailgun.InviteEmailParams{
				Name:             user.FirstName,
				Email:            user.Email,
				AcceptLink:       event.InviteURL,
				OrganizationName: org.Name,
				InvitorName:      invitorName,
			})

			if err != nil {
				gctx.Logger().Warnf("Unable to send invite email to %s %v", user.Email, err)
				return
			}

		}(gctx, agg, evt)
	}

	return nil
}
