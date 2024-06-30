package mailgun

import (
	"context"
	"fmt"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
)

type Email struct {
	Recipient string

	Subject string
	Body    string
	Timeout time.Duration
}

type EmailWithTemplate struct {
	Email

	Template  string
	Variables map[string]interface{}
}

func (*DefaultClient) SendEmail(gctx golly.Context, email Email) error {
	return nil
}
func (c *DefaultClient) SendEmailTemplate(gctx golly.Context, email EmailWithTemplate) error {
	sender := "TalentRadar <no-reply@notifications.talent-radar.io>"
	if s := gctx.Config().GetString("mailgun.sender"); s != "" {
		sender = s
	}

	subject := email.Subject
	if !golly.Env().IsProduction() {
		subject = fmt.Sprintf("TEST: %s", subject)
	}

	message := c.mailgun.NewMessage(sender, subject, email.Body, email.Recipient)
	if email.Template == "" {
		return errors.WrapGeneric(fmt.Errorf("template name must be provided"))
	}

	message.SetTemplate(email.Template)

	if email.Variables != nil {
		for key, val := range email.Variables {
			message.AddVariable(key, val)
		}

	}

	timeout := 5 * time.Second
	if email.Timeout.Seconds() > 0 {
		timeout = email.Timeout
	}

	ctx, cancel := context.WithTimeout(gctx.Context(), timeout)
	defer cancel()

	gctx.Logger().Debugf("Sending Email: %#v", email)

	_, _, err := c.mailgun.Send(ctx, message)
	return err
}
