package users

import (
	"time"

	"github.com/google/uuid"
)

type UserInvited struct {
	IdpInviteID string

	InvitedAt *time.Time
	InviterID uuid.UUID

	InviteURL string
}

type UserCreated struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID

	IdpID string

	Email     string
	FirstName string
	LastName  string
}

type UserUpdated struct {
	IdpID          string
	Email          string
	FirstName      string
	LastName       string
	ProfilePicture string
}
