package reviews

import (
	"fmt"
	"strings"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
)

const (
	feedbackGroupQuery = `
	SELECT 
		employee_id, 
		%s as collection_end_at,
		organization_id, 
		string_agg(id%s, ',') AS feedback_ids,
		SUM(
			CASE
			WHEN submitted_at IS NOT NULL THEN 1
			ELSE 0
			END
		) as submitted_count,
		COUNT(0) as sent_count
	FROM feedbacks
	WHERE organization_id = @organizationID
		AND employee_id IN @employeeIDs
	GROUP BY 1, 2, 3
	ORDER BY employee_id, collection_end_at
	LIMIT @limit OFFSET @offset
`
)

type GroupedFeedbackResults struct {
	EmployeeID      uuid.UUID
	CollectionEndAt time.Time
	OrganizationID  uuid.UUID

	TotalSent      int
	TotalSubmitted int

	FeedbackIDS uuid.UUIDs
}

func GroupedFeedback(gctx golly.Context, managerID uuid.UUID, limit, offset int) ([]GroupedFeedbackResults, error) {
	var results []GroupedFeedbackResults

	ident := identity.FromContext(gctx)

	var rawResults []struct {
		SubmittedCount  int
		SentCount       int
		EmployeeID      uuid.UUID
		CollectionEndAt string `gorm:"type:date"`
		OrganizationID  uuid.UUID
		FeedbackIDs     string
	}

	subordIDs, err := employees.
		Service(gctx).
		PluckEmployeeIDsByManagerID(gctx, managerID)

	if err != nil {
		return []GroupedFeedbackResults{}, err
	}

	gctx.Logger().Debugf("Loading for %#v", subordIDs)

	str := ""
	if !golly.Env().IsTest() {
		str = "::character varying"
	}

	dateParser := "TO_CHAR(collection_end_at, 'YYYY-MM-DD')"
	if golly.Env().IsTest() {
		dateParser = "strftime('%Y-%m-%d', collection_end_at)"
	}

	query := fmt.Sprintf(feedbackGroupQuery, dateParser, str)

	err = orm.
		DB(gctx).
		Raw(query, map[string]interface{}{
			"organizationID": ident.OrganizationID,
			"employeeIDs":    subordIDs,
			"limit":          limit,
			"offset":         offset,
			"timezero":       time.Time{},
		}).
		Scan(&rawResults).
		Error

	if err != nil {
		return results, err
	}

	for _, raw := range rawResults {

		t, err := time.Parse("2006-01-02", raw.CollectionEndAt)
		if err != nil {
			return results, err
		}

		feedbackIDs := strings.Split(raw.FeedbackIDs, ",")

		results = append(results, GroupedFeedbackResults{
			EmployeeID:      raw.EmployeeID,
			CollectionEndAt: t,
			OrganizationID:  raw.OrganizationID,
			TotalSubmitted:  raw.SubmittedCount,
			TotalSent:       raw.SentCount,
			FeedbackIDS: golly.Map(feedbackIDs, func(id string) uuid.UUID {
				return uuid.MustParse(id)
			}),
		})
	}

	return results, nil
}
