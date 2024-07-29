package main

import (
	"strings"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/role"
	"github.com/mitchrodrigues/talent-review-backend/app/initializers"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var commands = []*cobra.Command{
	{
		Use:  "convert-roles",
		Long: "send test feedback email",
		Run:  golly.Command(convertRoles),
	},
}

type OldRole struct {
	Title          string
	Level          int
	OrganizationID uuid.UUID
	EmployeeIDs    string
	Track          role.EmployeeType
}

const query = `
	SELECT DISTINCT ON (organization_id, title)
		organization_id, 
		title,
		MAX(level) OVER (PARTITION BY organization_id, title) as level,
		string_agg(id::character varying, ',') OVER (PARTITION BY organization_id, title) as employee_ids,
		FIRST_VALUE(type) OVER (PARTITION BY organization_id, title ORDER BY created_at) as track
	FROM employees
`

func convertRoles(gctx golly.Context, cmd *cobra.Command, args []string) error {
	var oldRoles []OldRole

	{
		err := orm.
			DB(gctx).
			Raw(query).
			Scan(&oldRoles).
			Error

		if err != nil {
			return err
		}
	}

	return orm.DB(gctx).Transaction(func(txn *gorm.DB) error {

		for _, oldRole := range oldRoles {
			newRole := role.Aggregate{}

			if oldRole.Title == "" {
				continue
			}

			err := eventsource.Call(gctx, &newRole, role.CreateEmployeeRole{
				Title:          oldRole.Title,
				Level:          oldRole.Level,
				OrganizationID: oldRole.OrganizationID,
				Track:          oldRole.Track,
			}, eventsource.Metadata{"cli": "true"})

			if err != nil {
				if err == role.ErrorRoleExists {
					continue
				}
				return err
			}

			ids := golly.Map(strings.Split(oldRole.EmployeeIDs, ","), func(s string) uuid.UUID {
				return uuid.MustParse(s)
			})

			var emps []employees.Employee

			orm.DB(gctx).
				Model(emps).
				Find(&emps, "id IN ?", ids)

			for _, emp := range emps {
				eventsource.Call(gctx, &emp.Aggregate, employee.Update{
					EmployeeRoleID: newRole.ID,
				}, eventsource.Metadata{"cli": "true"})
			}
		}

		return nil
	})
}

func main() {
	golly.Start(golly.GollyStartOptions{
		Preboots:     initializers.Preboots,
		Initializers: initializers.Initializers,
		CLICommands:  commands,
	})
}
