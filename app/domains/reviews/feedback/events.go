package feedback

import (
	"time"

	"github.com/google/uuid"
)

type Created struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID

	Email string
	Code  string

	CollectionEndAt time.Time
}
