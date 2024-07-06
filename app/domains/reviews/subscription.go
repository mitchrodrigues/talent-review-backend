package reviews

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/mailgun"
	"gorm.io/gorm"
)

func SendFeedbackEmail(gctx golly.Context, agg eventsource.Aggregate, evt eventsource.Event) error {
	switch event := evt.Data.(type) {
	case feedback.Created:
		go func(ctx golly.Context, agg eventsource.Aggregate, evt eventsource.Event) {
			gctx = orm.SetDBOnContext(ctx, orm.Connection().Session(&gorm.Session{NewDB: true}))

			feedback := agg.(*feedback.Aggregate)

			user, err := accounts.FindUserForContext(gctx)
			if err != nil {
				return
			}

			employee, err := employees.FindEmployeeByID(gctx, feedback.EmployeeID)
			if err != nil {
				return
			}

			err = mailgun.GetClient(gctx).SendFeedbackEmail(gctx, mailgun.FeedbackEmailParams{
				Name:            employee.Name,
				Email:           feedback.Email,
				CollectionEndAt: event.CollectionEndAt,
				FeedbackURL:     ctx.Config().GetString("app.frontend.url") + "/feedback/form/" + feedback.Code,
			})

			if err != nil {
				gctx.Logger().Warnf("Unable to send invite email to %s %v", user.Email, err)
				return
			}

		}(gctx, agg, evt)
	}

	return nil
}
