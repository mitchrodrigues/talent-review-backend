package mailgun

import (
	"fmt"
	"time"

	"github.com/golly-go/golly"
)

// feedbackURL
// name
// date

type FeedbackEmailParams struct {
	Name            string
	Email           string
	FeedbackURL     string
	CollectionEndAt time.Time
}

func (c *DefaultClient) SendFeedbackEmail(gctx golly.Context, params FeedbackEmailParams) error {
	return c.SendEmailTemplate(gctx, EmailWithTemplate{
		Email: Email{
			Recipient: params.Email,
			Subject:   fmt.Sprintf("Feedback Request for %s", params.Name),
		},
		Template: "feedback request",
		Variables: map[string]interface{}{
			"name":        params.Name,
			"email":       params.Email,
			"feedbackURL": params.FeedbackURL,
			"date":        params.CollectionEndAt.Format("01/02/2006"),
		},
	})
}
