package teams

import "github.com/google/uuid"

type TeamCreated struct {
	ID             uuid.UUID
	ManagerID      uuid.UUID
	OrganizationID uuid.UUID

	Name string
}

type TeamUpdated struct {
	Name      string
	ManagerID uuid.UUID
}
