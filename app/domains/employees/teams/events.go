package teams

import "github.com/google/uuid"

type Created struct {
	ID             uuid.UUID
	LeadID         *uuid.UUID `json:"omitempty"`
	OrganizationID uuid.UUID
	Name           string
}

type Updated struct {
	Name   string
	LeadID *uuid.UUID `json:"teamID,omitempty"`
}
