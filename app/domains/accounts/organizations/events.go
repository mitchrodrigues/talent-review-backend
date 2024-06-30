package organizations

import (
	"github.com/google/uuid"
)

type OrganizationCreated struct {
	ID       uuid.UUID `json:"id"`
	IdpID    string    `json:"idpID"`
	Name     string    `json:"name"`
	PlanName string    `json:"planName"`
}
