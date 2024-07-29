package employee

import (
	"time"

	"github.com/google/uuid"
)

type Created struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	OrganizationID uuid.UUID `json:"organizationID"`
}

type Updated struct {
	Name       string             `json:"name"`
	Email      string             `json:"email"`
	Level      int                `json:"level"`
	WorkerType EmployeeWorkerType `json:"workerType"`
}

type UserUpdated struct {
	UserID uuid.UUID `json:"userID"`
}

type TeamUpdated struct {
	TeamID *uuid.UUID `json:"teamID"`
}

type TitleUpdated struct {
	Title string `json:"teamID"`
}

type ManagerUpdated struct {
	ManagerID *uuid.UUID `json:"managerID"`
}

type PersonalDetailsUpdated struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type WorkerTypeUpdated struct {
	WorkerType EmployeeWorkerType `json:"workerType"`
}

type RoleUpdated struct {
	EmployeeRoleID uuid.UUID
}

type Terminate struct {
	TerminatedAt time.Time
}
