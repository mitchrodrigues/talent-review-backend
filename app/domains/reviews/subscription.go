package reviews

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/tara"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/mailgun"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/wsyiwig"
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

			employee, err := employees.Service(gctx).FindEmployeeByID(gctx, feedback.EmployeeID)
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

func UpdateFeedbackSummarySubscription(gctx golly.Context, agg eventsource.Aggregate, evt eventsource.Event) error {
	switch evt.Data.(type) {
	case feedback.Submitted:
		fb := agg.(*feedback.Aggregate)

		go func(gctx golly.Context, fb *feedback.Aggregate) {
			if err := UpdateFeedbackSummary(gctx, fb); err != nil {
				gctx.Logger().Warnf("connot generate summary for feedback %s %v", fb.ID.String(), err)
			}
		}(gctx, fb)

		return nil
	}
	return nil
}

func UpdateFeedbackSummary(gctx golly.Context, fb *feedback.Aggregate) error {

	details, err := Service(gctx).FindFeedbackDetailsByFeedbackID_Unsafe(gctx, fb.ID)
	if err != nil {
		gctx.Logger().Warnf("cannot find details for feedback %s %v", fb.ID.String(), err)
		return err
	}

	strengths, _ := wsyiwig.ExtractTextFromJSON(details.Strengths)
	opportunities, _ := wsyiwig.ExtractTextFromJSON(details.Opportunities)
	additional, _ := wsyiwig.ExtractTextFromJSON(details.Additional)

	prompt := tara.NewSummaryFeedbackPrompt(tara.SummarizeFeedbackInput{
		Strengths:          strengths,
		Opportunities:      opportunities,
		AdditionalComments: additional,
	})

	err = tara.Generate(gctx, prompt)
	if err != nil {
		return errors.WrapGeneric(err)
	}

	itemsPrompt := tara.NewFollowUpItemsPrompt()
	itemsPrompt.AddPreviousPrompts(prompt)

	err = tara.Generate(gctx, itemsPrompt)
	if err != nil {
		return errors.WrapGeneric(err)
	}

	err = eventsource.Call(gctx, fb, feedback.CreateSummary{
		Summary:     prompt.Summary,
		ActionItems: itemsPrompt.FollowUpItems.Values(),
	}, eventsource.Metadata{})

	return errors.WrapGeneric(err)
}
