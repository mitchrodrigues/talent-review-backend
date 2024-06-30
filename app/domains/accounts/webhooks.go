package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts/users"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
)

type WebhookController struct{}

func (controller WebhookController) Routes(router *golly.Route) {
	router.Post("/workos/users", controller.workosUserEvent)
}

func (cntr WebhookController) workosUserEvent(wctx golly.WebContext) {
	var event workos.UserWebhookEvent

	if err := json.Unmarshal(wctx.RequestBody(), &event); err != nil {
		wctx.RenderStatus(http.StatusUnprocessableEntity)
		return
	}

	user, err := FindUserByEmail(wctx.Context, event.Data.Email)
	if err == nil {

		err := eventsource.Call(wctx.Context, &user, users.EditUser{
			ProfilePicture: event.Data.ProfilePictureURL,
			FirstName:      event.Data.FirstName,
			LastName:       event.Data.LastName,
			IdpID:          event.Data.ID,
		}, eventsource.Metadata{})

		if err != nil {
			wctx.RenderStatus(http.StatusUnprocessableEntity)
			return
		}
	}

	wctx.RenderStatus(http.StatusOK)
}
