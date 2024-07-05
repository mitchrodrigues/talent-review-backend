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
	results := []Feedback{}
	db := orm.DB(gctx)

	emps, err := employees.FindEmployeesByIDS(gctx, input.EmployeeIDs)
	if err != nil {
		return nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		tCtx := orm.SetDBOnContext(gctx.Dup(), tx)

		for _, employee := range emps {
			gctx.Logger().Debugf("Starting Process Of Bulk Feedback: %#v", employee)

			if aggs, err := processAdditionalEmails(tCtx, input, employee, metadata); err != nil {
				return err
			} else {
				results = append(results, aggs...)
			}

			if !input.IncludeTeam || employee.TeamID == nil {
				continue
			}

			if aggs, err := processTeamMembers(tCtx, input, employee, metadata); err != nil {
				return err
			} else {
				results = append(results, aggs...)
			}
		}
		return nil
	})

	return results, err
}

func processAdditionalEmails(gctx golly.Context,
	input CreateBulkFeedbackInput,
	employee employees.Employee,
	metadata eventsource.Metadata,
) ([]Feedback, error) {
	ident := identity.FromContext(gctx)

	var results []Feedback

	gctx.Logger().Debugf("Starting Processing Additional Emails: %#v", input.AdditionalEmails)

	for _, email := range input.AdditionalEmails {
		if strings.EqualFold(email, employee.Email) {
			continue
		}

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
	return results, nil
}

func processTeamMembers(gctx golly.Context,
	input CreateBulkFeedbackInput,
	employee employees.Employee,
	metadata eventsource.Metadata,
) ([]Feedback, error) {
	gctx.Logger().Debugf("Starting Processing TeamMembers Emails: %#v", input.IncludeTeam)

	ident := identity.FromContext(gctx)

	results := []Feedback{}

	teamMates, err := getTeamMates(gctx, employee.TeamID)
	if err != nil {
		return results, err
	}

	for _, teamMate := range teamMates {
		record := Feedback{}
		err := eventsource.Call(gctx, &record.Aggregate, feedback.Create{
			CollectionEndAt: input.CollectionEndAt,
			EmployeeID:      teamMate.ID,
			OrganizationID:  ident.OrganizationID,
			Email:           teamMate.Email,
		}, metadata)

		if err != nil {
			return results, err
		}

		results = append(results, record)
	}

	return results, nil
}

func getTeamMates(gctx golly.Context, teamID *uuid.UUID) ([]employees.Employee, error) {
	key := fmt.Sprintf("employees::team::%s", teamID.String())

	return golly.LoadData(gctx, key, func(golly.Context) ([]employees.Employee, error) {
		return employees.FindEmployeesForTeam(gctx, *teamID)
	})
}
