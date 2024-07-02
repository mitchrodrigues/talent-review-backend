package cycle

import (
	"time"

	"github.com/google/uuid"
)

type CycleCreated struct {
	ID             uuid.UUID
	OwnerID        uuid.UUID
	OrganizationID uuid.UUID
	Type           string

	StartAt time.Time
	EndAt   time.Time
}
