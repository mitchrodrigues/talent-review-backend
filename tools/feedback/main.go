package main

import (
	"github.com/golly-go/golly"
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
}

func main() {
	golly.Start(golly.GollyStartOptions{
		Preboots:     initializers.Preboots,
		Initializers: initializers.Initializers,
		CLICommands:  commands,
	})
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
			"feedbackURL": "http://localhost:9009/feedback/1234-1234",
		},
	})
}
