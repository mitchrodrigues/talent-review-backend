package feedback

import (
	"time"

	"github.com/google/uuid"
)

type Created struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
	OwnerID        uuid.UUID

	Email string
	Code  string

	CollectionEndAt time.Time
}

type Submitted struct{}

type DetailsCreated struct {
	ID             uuid.UUID
	FeedbackID     uuid.UUID
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

type DetailsUpdated struct {
	Strenghts     string `json:"-"`
	Opportunities string `json:"-"`
	Additional    string `json:"-"`
	EnoughData    bool   `json:"-"`

	Rating int `json:"-"`
}

type SummaryCreated struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationID"`
	EmployeeID     uuid.UUID `json:"employeeID"`
	FeedbackID     uuid.UUID `json:"feedbackID"`
}

type SummaryUpdated struct {
	Summary     string `json:"-"`
	ActionItems string `json:"-"`
}
