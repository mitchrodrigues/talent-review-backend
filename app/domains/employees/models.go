package employees

import (
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/role"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/teams"
)

type Employee struct {
	employee.Aggregate

	Team Team

	Role EmployeeRole `gorm:"foreignKey:EmployeeRoleID"`
}

func (Employee) TableName() string { return "employees" }

type Team struct {
	teams.Aggregate

	Manager   *Employee
	Employees []Employee
}

func (Team) TableName() string { return "teams" }

type EmployeeRole struct {
	role.Aggregate

	Employees []Employee
}

func (EmployeeRole) TableName() string { return "employee_roles" }
