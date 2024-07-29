package employees

import (
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/teams"
)

func NewTestEmployee(id, organizationID uuid.UUID, email string, userID *uuid.UUID) Employee {
	return Employee{
		Aggregate: employee.Aggregate{
			ModelUUID: orm.ModelUUID{
				ID: id,
			},
			Email:          email,
			OrganizationID: organizationID,
			UserID:         userID,
		},
	}
}

func NewTestEmployeeWithTeam(id, organizationID uuid.UUID, email string, teamID uuid.UUID) Employee {
	employee := NewTestEmployee(id, organizationID, email, nil)
	employee.TeamID = &teamID

	return employee
}

func NewTestTeam(id, organizationID uuid.UUID, leadID *uuid.UUID) Team {
	return Team{
		Aggregate: teams.Aggregate{
			ModelUUID: orm.ModelUUID{
				ID: id,
			},
			LeadID:         leadID,
			OrganizationID: organizationID,
		},
	}
}
