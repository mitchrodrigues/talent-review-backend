package employees

import (
	"encoding/json"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
)

func UpdateEmployeeUser(gctx golly.Context, agg eventsource.Aggregate, event eventsource.Event) error {
	// I hate this but not sure how to cast this better
	b, _ := json.Marshal(event.Data)
	var data map[string]interface{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	email, err := helpers.GetString(data, "Email")
	if err != nil {
		return nil
	}

	orgID, err := helpers.ExtractAndParseUUID(data, "OrganizationID")
	if err != nil || orgID == uuid.Nil {
		return nil
	}

	readModel, err := Service(gctx).FindEmployeeByEmailAndOrganizationID(gctx, email, orgID)
	if readModel.ID == uuid.Nil || err != nil {
		return err
	}

	return eventsource.Handler(gctx).Call(gctx, &readModel.Aggregate, employee.UpdateUser{
		UserID: uuid.MustParse(event.AggregateID),
	}, event.Metadata)
}
