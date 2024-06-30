package employees

import (
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/teams"
)

type Employee struct {
	employee.Aggregate

	Team Team
}

func (Employee) TableName() string { return "employees" }

type Team struct {
	teams.Aggregate

	Manager   *Employee
	Employees []Employee
}

func (Team) TableName() string { return "teams" }
