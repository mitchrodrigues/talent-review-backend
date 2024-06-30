package identity

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/passport"
	"github.com/google/uuid"
)

type Identity struct {
	UID            uuid.UUID
	OrganizationID uuid.UUID
}

var _ passport.Identity = Identity{}

func (Identity) Valid() error           { return nil }
func (ident Identity) IsLoggedIn() bool { return ident.UID != uuid.Nil }
func (ident Identity) UserID() string   { return ident.UID.String() }

func Cast(ident passport.Identity) Identity {
	if ident, ok := ident.(Identity); ok {
		return ident
	}
	return Identity{}
}

func FromContext(gctx golly.Context) Identity {
	if ident, found := passport.FromContext(gctx); found {
		return ident.(Identity)
	}
	return Identity{}
}
