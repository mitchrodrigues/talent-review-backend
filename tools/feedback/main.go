package main

import (
	"encoding/json"
	"fmt"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews"
	"github.com/mitchrodrigues/talent-review-backend/app/initializers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/mailgun"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var commands = []*cobra.Command{
	{
		Use:  "test-email",
		Long: "send test feedback email",
		Args: cobra.MinimumNArgs(1),
		Run:  golly.Command(testEmail),
	},
	{
		Use:  "update-summary [feedbackID]",
		Long: "update tara summary for a feedback",
		Args: cobra.MinimumNArgs(1),
		Run:  golly.Command(submitFeedback),
	},
	{
		Use:  "check [managerID]",
		Long: "update tara summary for a feedback",
		Args: cobra.MinimumNArgs(1),
		Run:  golly.Command(checkFeedbacks),
	},
}

func main() {
	golly.Start(golly.GollyStartOptions{
		Preboots:     initializers.Preboots,
		Initializers: initializers.Initializers,
		CLICommands:  commands,
	})
}

func checkFeedbacks(gctx golly.Context, cmd *cobra.Command, args []string) error {
	logger := gctx.Logger()

	logger.Logger.SetLevel(logrus.DebugLevel)
	gctx.SetLogger(logger)

	db := orm.DB(gctx).Session(&gorm.Session{
		NewDB:  true,
		Logger: orm.NewLogger("checkFeedbacks", false),
	})

	gctx = orm.SetDBOnContext(gctx, db)

	id := uuid.MustParse(args[0])

	results, err := reviews.GroupedFeedback(gctx, id, 50, 0)

	b, _ := json.MarshalIndent(results, "", "\t")

	fmt.Printf("GroupedFeedback: \n%s\n", string(b))

	return err
}

func submitFeedback(gctx golly.Context, cmd *cobra.Command, args []string) error {
	fb, err := reviews.FeedbackService(gctx).FindByID_Unsafe(gctx, uuid.MustParse(args[0]))
	if err != nil {
		return err
	}

	if fb.ID == uuid.Nil {
		return fmt.Errorf("no such feedback %s", args[0])
	}

	return reviews.UpdateFeedbackSummary(gctx, &fb.Aggregate)
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
