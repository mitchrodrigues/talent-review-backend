package reviews

import (
	"fmt"
	"strings"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
)

type CreateBulkFeedbackInput struct {
	EmployeeIDs      []uuid.UUID
	AdditionalEmails []string
	IncludeTeam      bool
	IncludeDirects   bool
	CollectionEndAt  time.Time
}

func CreateBulkFeedback(gctx golly.Context, input CreateBulkFeedbackInput, metadata eventsource.Metadata) ([]Feedback, error) {
	ident := identity.FromContext(gctx)

	results := []Feedback{}

	manager, err := employees.Service(gctx).FindEmployeeByUserID(gctx, ident.UID)
	if err != nil {
		return []Feedback{}, errors.WrapGeneric(fmt.Errorf("you are not a manager of any team"))
	}

	emps, err := employees.Service(gctx).FindEmployeesByManagerAndIDS(gctx, manager.ID, input.EmployeeIDs...)
	if err != nil {
		return nil, err
	}

	if len(emps) == 0 {
		return nil, errors.WrapGeneric(fmt.Errorf("you do not have any employees matching the criteria"))
	}

	for _, employee := range emps {
		emails := input.AdditionalEmails

		if input.IncludeTeam {
			teamMates, err := getTeamMates(gctx, employee.TeamID)
			if err != nil {
				return results, err
			}

			emails = append(emails, golly.Map(teamMates, func(employee employees.Employee) string {
				return employee.Email
			})...)
		}

		if input.IncludeDirects {
			emps, _ := employees.Service(gctx).FindEmployeesByManagerID(gctx, employee.ID)
			if len(emps) > 0 {
				emails = append(emails, golly.Map(emps, func(e employees.Employee) string {
					return e.Email
				})...)
			}
		}

		emails = golly.Unique(emails)

		for _, email := range emails {
			if strings.EqualFold(email, employee.Email) {
				continue
			}

			gctx.Logger().Debugf("Starting Process Of Bulk Feedback: %#v", employee)

			record := Feedback{}

			err := eventsource.Call(gctx, &record.Aggregate, feedback.Create{
				CollectionEndAt: input.CollectionEndAt,
				EmployeeID:      employee.ID,
				OrganizationID:  ident.OrganizationID,
				Email:           email,
			}, metadata)

			if err != nil {
				return results, err
			}

			results = append(results, record)
		}

	}
	return results, err
}

func getTeamMates(gctx golly.Context, teamID *uuid.UUID) ([]employees.Employee, error) {
	if teamID == nil {
		return []employees.Employee{}, nil
	}

	key := fmt.Sprintf("employees::team::%s", teamID.String())
	return golly.LoadData(gctx, key, func(golly.Context) ([]employees.Employee, error) {
		return employees.Service(gctx).FindEmployeesForTeam(gctx, *teamID)
	})
}
