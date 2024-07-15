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
		MAX(collection_end_at) as collection_end_at, 
		organization_id, 
		string_agg(id%s, ',') AS feedback_ids,
		SUM(
			CASE
			WHEN submitted_at IS NOT NULL THEN 1
			ELSE 0
			END
		) as submitted_count
	FROM feedbacks
	WHERE organization_id = @organizationID
		AND employee_id IN @employeeIDs
	GROUP BY employee_id, CAST(collection_end_at AS DATE), organization_id
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

	FeedbackIDS []uuid.UUID
}

func GroupedFeedback(gctx golly.Context, managerID uuid.UUID, limit, offset int) ([]GroupedFeedbackResults, error) {
	var results []GroupedFeedbackResults

	ident := identity.FromContext(gctx)

	var rawResults []struct {
		SubmittedCount  int
		EmployeeID      uuid.UUID
		CollectionEndAt time.Time
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

	query := fmt.Sprintf(feedbackGroupQuery, str)

	err = orm.
		DB(gctx).
		Raw(query, map[string]interface{}{
			"organizationID": ident.OrganizationID,
			"employeeIDs":    subordIDs,
			"limit":          limit,
			"offset":         offset,
		}).
		Scan(&rawResults).
		Error

	if err != nil {
		return results, nil
	}

	for _, raw := range rawResults {
		feedbackIDs := golly.Map(strings.Split(raw.FeedbackIDs, ","), func(s string) uuid.UUID {
			return uuid.MustParse(s)
		})

		if err != nil {
			return results, err
		}

		results = append(results, GroupedFeedbackResults{
			EmployeeID:      raw.EmployeeID,
			CollectionEndAt: raw.CollectionEndAt,
			OrganizationID:  raw.OrganizationID,
			TotalSubmitted:  raw.SubmittedCount,
			TotalSent:       len(feedbackIDs),
			FeedbackIDS:     feedbackIDs,
		})
	}

	return results, nil
}
