package main

import (
	"fmt"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/mitchrodrigues/talent-review-backend/app/initializers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/mailgun"
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	{
		Use:  "test-email",
		Long: "send test feedback email",
		Args: cobra.MinimumNArgs(1),
		Run:  golly.Command(testEmail),
	},
	{
		Use:  "submit-feedback [feedbackID]",
		Long: "submit a feedback",
		Args: cobra.MinimumNArgs(1),
		Run:  golly.Command(submitFeedback),
	},
}

func main() {
	golly.Start(golly.GollyStartOptions{
		Preboots:     initializers.Preboots,
		Initializers: initializers.Initializers,
		CLICommands:  commands,
	})
}

func submitFeedback(gctx golly.Context, cmd *cobra.Command, args []string) error {
	fb, err := reviews.Service(gctx).FindFeedbackByID_Unsafe(gctx, uuid.MustParse(args[0]))
	if err != nil {
		return err
	}

	if fb.ID == uuid.Nil {
		return fmt.Errorf("no such feedback %s", args[0])
	}

	_ = eventsource.Call(gctx, &fb.Aggregate, feedback.Submit{}, eventsource.Metadata{})

	return nil
}

func testEmail(gctx golly.Context, cmd *cobra.Command, args []string) error {
	mg := mailgun.NewDefaultClient(gctx)

	return mg.SendEmailTemplate(gctx, mailgun.EmailWithTemplate{
		Template: "feedback request",
		Email: mailgun.Email{
			Recipient: args[0],
			Subject:   "Feedback Email",
		},
		Variables: map[string]interface{}{
			"name":        "Test Employee",
			"feedbackURL": fmt.Sprintf("%s/feedback/1234-1234", gctx.Config().GetString("app.frontend.url")),
		},
	})
}
