package employee

import "github.com/google/uuid"

type Created struct {
	ID             uuid.UUID          `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	OrganizationID uuid.UUID          `json:"organizationID"`
	Level          int                `json:"level"`
	Type           EmployeeType       `json:"type"`
	WorkerType     EmployeeWorkerType `json:"workerType"`
	UserID         *uuid.UUID         `json:"userID"`
}

type UserUpdated struct {
	UserID uuid.UUID `json:"userID"`
}

type TeamUpdated struct {
	TeamID uuid.UUID `json:"teamID"`
}

type TitleUpdated struct {
	Title string `json:"teamID"`
}
