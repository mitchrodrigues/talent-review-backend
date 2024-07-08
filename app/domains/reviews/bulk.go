package reviews

import (
	"fmt"
	"strings"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"gorm.io/gorm"
)

type CreateBulkFeedbackInput struct {
	EmployeeIDs      []uuid.UUID
	AdditionalEmails []string
	IncludeTeam      bool
	CollectionEndAt  time.Time
}

func CreateBulkFeedback(gctx golly.Context, input CreateBulkFeedbackInput, metadata eventsource.Metadata) ([]Feedback, error) {
	ident := identity.FromContext(gctx)

	results := []Feedback{}
	db := orm.DB(gctx)

	emps, err := employees.Service(gctx).FindEmployeesByIDS(gctx, input.EmployeeIDs)
	if err != nil {
		return nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		tCtx := orm.SetDBOnContext(gctx.Dup(), tx)

		for _, employee := range emps {
			emails := input.AdditionalEmails

			if input.IncludeTeam {
				teamMates, err := getTeamMates(gctx, employee.TeamID)
				if err != nil {
					return err
				}

				emails = append(emails, golly.Map(teamMates, func(employee employees.Employee) string {
					return employee.Email
				})...)
			}

			emails = golly.Unique(emails)

			for _, email := range emails {
				if strings.EqualFold(email, employee.Email) {
					continue
				}

				gctx.Logger().Debugf("Starting Process Of Bulk Feedback: %#v", employee)

				record := Feedback{}

				err := eventsource.Call(tCtx, &record.Aggregate, feedback.Create{
					CollectionEndAt: input.CollectionEndAt,
					EmployeeID:      employee.ID,
					OrganizationID:  ident.OrganizationID,
					Email:           email,
				}, metadata)

				if err != nil {
					return err
				}

				results = append(results, record)
			}

		}
		return nil
	})

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
