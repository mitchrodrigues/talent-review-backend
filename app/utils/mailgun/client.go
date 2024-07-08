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
	SendFeedbackEmail(golly.Context, FeedbackEmailParams) error
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
func (m *MockEmailClient) SendEmail(gctx golly.Context, email Email) error {
	args := m.Called(gctx, email)
	return args.Error(0)
}

// SendEmailTemplate mocks the SendEmailTemplate method.
func (m *MockEmailClient) SendEmailTemplate(gctx golly.Context, emailTemplate EmailWithTemplate) error {
	args := m.Called(gctx, emailTemplate)
	return args.Error(0)
}

// SendEmailTemplate mocks the SendEmailTemplate method.
func (m *MockEmailClient) SendInviteEmail(gctx golly.Context, params InviteEmailParams) error {
	args := m.Called(gctx, params)
	return args.Error(0)
}

// SendEmailTemplate mocks the SendEmailTemplate method.
func (m *MockEmailClient) SendFeedbackEmail(gctx golly.Context, params FeedbackEmailParams) error {
	args := m.Called(gctx, params)
	return args.Error(0)
}

func GetClient(ctx golly.Context) Client {
	if client, found := ctx.Get(contextKey); found {
		return client.(Client)
	}

	if golly.Env().IsTest() {
		return &MockEmailClient{}
	}

	return NewDefaultClient(ctx)
}
