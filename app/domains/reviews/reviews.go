package reviews

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
)

func Initializer(app golly.Application) error {
	InitGraphQL()

	eventsource.Subscribe("feedback.Aggregate", "feedback.Created", SendFeedbackEmail)
	eventsource.Subscribe("feedback.Aggregate", "feedback.Submitted", UpdateFeedbackSummarySubscription)

	return nil
}
