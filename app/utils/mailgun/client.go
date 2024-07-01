package mailgun

import (
	"github.com/golly-go/golly"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/mock"
)

const (
	contextKey golly.ContextKeyT = "mailgunClient"
)

type Client interface {
	SendEmail(golly.Context, Email) error
	SendEmailTemplate(golly.Context, EmailWithTemplate) error
	SendInviteEmail(golly.Context, InviteEmailParams) error
}
type DefaultClient struct {
	mailgun *mailgun.MailgunImpl
}

func NewDefaultClient(gctx golly.Context) Client {
	return &DefaultClient{
		mailgun: mailgun.NewMailgun(
			gctx.Config().GetString("mailgun.domain"),
			gctx.Config().GetString("mailgun.key"),
		),
	}
}

// MockEmailClient is a mock implementation of the Client interface using testify/mock.
type MockEmailClient struct {
	mock.Mock
}

// SendEmail mocks the SendEmail method.
func (m *MockEmailClient) SendEmail(email Email) error {
	args := m.Called(email)
	return args.Error(0)
}

// SendEmailTemplate mocks the SendEmailTemplate method.
func (m *MockEmailClient) SendEmailTemplate(emailTemplate EmailWithTemplate) error {
	args := m.Called(emailTemplate)
	return args.Error(0)
}

func GetClient(ctx golly.Context) Client {
	if client, found := ctx.Get(contextKey); found {
		return client.(Client)
	}

	return NewDefaultClient(ctx)
}
