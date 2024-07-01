package workos

import (
	"net/url"

	"github.com/golly-go/golly"
	"github.com/google/uuid"
	"github.com/workos/workos-go/v4/pkg/organizations"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
)

const (
	contextKey golly.ContextKeyT = "workos"
)

type IDPObject interface {
	RecordIdpID() string
	RecordID() uuid.UUID
}

type WorkosClient interface {
	CreateUser(golly.Context, CreateUserInput) (string, error)
	CreateOrganization(golly.Context, string, ...string) (string, error)
	InviteUser(golly.Context, string, string, string) (string, string, error)
	JWKSURL() (*url.URL, error)
}

type DefaultClient struct {
	ClientID string
}

func (d DefaultClient) JWKSURL() (*url.URL, error) {
	return usermanagement.GetJWKSURL(d.ClientID)
}

type UserWebhookEvent struct {
	ID    string              `json:"id"`
	Event string              `json:"event"`
	Data  usermanagement.User `json:"data"`
}

func Initializer(app golly.Application) error {
	key := app.Config.GetString("workos.api.key")

	organizations.SetAPIKey(key)
	usermanagement.SetAPIKey(key)

	return nil
}

func Client(ctx golly.Context) WorkosClient {
	if client, found := ctx.Get(contextKey); found {
		return client.(WorkosClient)
	}

	if golly.Env().IsTest() {
		return &MockClient{}
	}

	clientID := ctx.Config().GetString("workos.client.id")

	return DefaultClient{clientID}
}
