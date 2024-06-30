package users

import (
	"time"

	"github.com/google/uuid"
)

type UserCreated struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID

	IdpID       string
	IdpInviteID string

	Email     string
	FirstName string
	LastName  string

	InvitedAt *time.Time
	InviterID uuid.UUID
}

type UserUpdated struct {
	IdpID          string
	Email          string
	FirstName      string
	LastName       string
	ProfilePicture string
}
