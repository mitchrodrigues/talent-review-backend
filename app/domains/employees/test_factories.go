package employees

import (
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/teams"
)

func NewTestEmployee(id, organizationID uuid.UUID, userID *uuid.UUID) Employee {
	return Employee{
		Aggregate: employee.Aggregate{
			ModelUUID: orm.ModelUUID{
				ID: id,
			},
			OrganizationID: organizationID,
			UserID:         userID,
		},
	}
}

func NewTestTeam(id, organizationID, managerID uuid.UUID) Team {
	return Team{
		Aggregate: teams.Aggregate{
			ModelUUID: orm.ModelUUID{
				ID: id,
			},
			ManagerID:      managerID,
			OrganizationID: organizationID,
		},
	}
}
